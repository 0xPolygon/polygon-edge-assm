package localstorage

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// pin the relese to v0.5.1 - before BLS breaking change
const releaseURL string = "https://api.github.com/repos/0xPolygon/polygon-edge/releases/77894422"

type Adapter struct {
}

type Asset struct {
	DownloadURL string `json:"browser_download_url"`
}

type PolygonEdgeRelase struct {
	Assets           []Asset `json:"assets"`
	realDownloadLink string
}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a Adapter) GetEdge() error {
	edge := &PolygonEdgeRelase{}

	// get the info about releases
	resp, err := http.Get(releaseURL)
	if err != nil {
		return fmt.Errorf("could not get new polygon-edge release err=%w", err)
	}
	defer resp.Body.Close()

	// read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf(
			"could not read response body of "+releaseURL+" err=%w",
			err)
	}

	// unmarshal response
	if err := json.Unmarshal(body, edge); err != nil {
		return fmt.Errorf("could not unmarshal response body of "+
			releaseURL+"err=%w", err)
	}

	// find release for our platform - linux_amd64
	for _, link := range edge.Assets {
		if strings.Contains(link.DownloadURL, "linux_amd64") {
			log.Println("Downloading: ", link.DownloadURL)
			edge.realDownloadLink = link.DownloadURL
		}
	}

	// get edge release
	edgeResp, err := http.Get(edge.realDownloadLink)
	if err != nil {
		return fmt.Errorf("could not download release err=%w", err)
	}
	defer edgeResp.Body.Close()

	// read tar stream
	uncompressedStream, err := gzip.NewReader(edgeResp.Body)
	if err != nil {
		return fmt.Errorf("could not uncompress tar err=%w", err)
	}

	// untar release package
	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		//nolint:errorlint
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("ExtractTarGz: Next() failed err=%w", err)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir("/tmp/"+header.Name, 0755); err != nil {
				return fmt.Errorf("ExtractTarGz: Mkdir() failed err=%w", err)
			}
		case tar.TypeReg:
			//nolint:gosec
			outFile, err := os.Create("/tmp/" + header.Name)
			if err != nil {
				return fmt.Errorf("ExtractTarGz: Create() failed: %w", err)
			}
			//nolint:gosec
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("ExtractTarGz: Copy() failed: %w", err)
			}

			err = outFile.Chmod(fs.ModePerm)
			if err != nil {
				log.Println("could not set permissions on polygon-edge binary")
			}

			_ = outFile.Close()

		default:
			return fmt.Errorf(
				"ExtractTarGz: uknown type: %s in %s",
				string(header.Typeflag),
				header.Name)
		}
	}

	_ = os.Remove("/tmp/LICENSE")
	_ = os.Remove("/tmp/README.md")

	return nil
}

func (a Adapter) RunGenesisCmd(args []string) error {
	var out bytes.Buffer

	cmd := exec.Command("/tmp/polygon-edge", args...)
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not run genesis command err=%w", err)
	}

	// for debugging
	fmt.Printf("genesis output: %s\n", out.String())

	return nil
}
