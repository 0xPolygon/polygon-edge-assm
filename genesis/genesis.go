package genesis

import (
	"Trapesys/polygon-edge-assm/aws"
	"Trapesys/polygon-edge-assm/types"
	"fmt"
	"os"
	"os/exec"
)

func GenerateAndStore(nodes *types.Nodes) error {
	// first part of genesis command
	genCmd := []string{"genesis", "--consensus", "ibft", "--dir", "/tmp"}
	
	// add validators and keys
	for _, v := range nodes.Node {
		genCmd = append(genCmd, "--ibft-validator="+v.ValidatorKey)
		genCmd = append(genCmd, fmt.Sprintf("--bootnode=/ip4/%s/tcp/1478/p2p/%s", v.IP, v.NetworkID))
	}

	genCmd = append(genCmd, "--premine=0x228466F2C715CbEC05dEAbfAc040ce3619d7CF0B:1000000000000000000000")

	// remove temp file if exists
	os.Remove("/tmp/genesis.json")

	_, err := exec.Command("polygon-edge",genCmd...).Output()
	if err != nil {
		return fmt.Errorf("failed to generate genesis.json: %w", err)
	}

	fmt.Println(genCmd)

	if err := aws.StoreGenesis("/tmp/genesis.json"); err != nil {
		return fmt.Errorf("failed to store genesis.json: %w",err)
	}

	return nil
}