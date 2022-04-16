package genesis

import (
	"Trapesys/polygon-edge-assm/aws"
	"Trapesys/polygon-edge-assm/types"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var genesisPath = "/tmp/genesis.json"

func GenerateAndStore(nodes *types.Nodes) error {
	// first part of genesis command
	genCmd := []string{"genesis", "--consensus", "ibft", "--dir", genesisPath}
	
	// add validators and keys
	for _, v := range nodes.Node {
		genCmd = append(genCmd, "--ibft-validator="+v.ValidatorKey)
		genCmd = append(genCmd, fmt.Sprintf("--bootnode=/ip4/%s/tcp/1478/p2p/%s", strings.TrimSpace(v.IP), strings.TrimSpace(v.NetworkID)))
	}

	genCmd = append(genCmd, "--premine=0x228466F2C715CbEC05dEAbfAc040ce3619d7CF0B:1000000000000000000000")

	// remove temp file if exists
	os.Remove(genesisPath)

	cmd := exec.Command("polygon-edge",genCmd...)

	logWriter, err := os.Create("/var/log/edge-controler.log")
	if err != nil {
		return fmt.Errorf("could not setup log file writer, %w",err)
	}

	cmd.Stdout = logWriter
	cmd.Stderr = logWriter

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to generate genesis.json: %w", err)
	}

	if err := aws.StoreGenesis(genesisPath); err != nil {
		return fmt.Errorf("failed to store genesis.json: %w",err)
	}

	return nil
}