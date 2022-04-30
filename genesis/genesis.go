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

type Config struct {
	ChainName string
	PoS bool
	EpochSize string
	Premine string
	ChainID string
	BlockGasLimit string
	MaxValidatorCount string
	MinValidatorCount string

	LogFile string
}

var GenConfig = &Config{}

func GenerateAndStore(nodes *types.Nodes) error {
	// first part of genesis command
	genCmd := []string{"genesis", "--consensus", "ibft", "--dir", genesisPath}

	// Name
	if GenConfig.ChainName != "" {
		genCmd = append(genCmd, "--name", GenConfig.ChainName)
	}

	// PoS params
	if GenConfig.PoS {
		genCmd = append(genCmd, "--pos")
	}

	// Epoch Size
	if GenConfig.EpochSize != "" {
		genCmd = append(genCmd, "--epoch-size", GenConfig.EpochSize)
	}

	// Chain ID
	if GenConfig.ChainID != "" {
		genCmd = append(genCmd, "--chain-id", GenConfig.ChainID)
	}

	// Block Gas Limit
	if GenConfig.BlockGasLimit != "" {
		genCmd = append(genCmd, "--block-gas-limit", GenConfig.BlockGasLimit)
	}

	// Max validator count
	if GenConfig.PoS && GenConfig.MaxValidatorCount != "" {
		genCmd = append(genCmd, "--max-validator-count", GenConfig.MaxValidatorCount)
	}

	// Min validator count
	if GenConfig.PoS && GenConfig.MinValidatorCount != "" {
		genCmd = append(genCmd, "--min-validator-count", GenConfig.MinValidatorCount)
	}

	
	// add validators and keys
	for _, v := range nodes.Node {
		genCmd = append(genCmd, "--ibft-validator="+v.ValidatorKey)
		genCmd = append(genCmd, fmt.Sprintf("--bootnode=/ip4/%s/tcp/1478/p2p/%s", strings.TrimSpace(v.IP), strings.TrimSpace(v.NetworkID)))
	}

	// Premine
	for _,premine := range strings.Split(GenConfig.Premine, ",") {
		genCmd = append(genCmd, fmt.Sprintf("--premine=%s",premine))
	}

	// remove temp file if exists
	os.Remove(genesisPath)

	cmd := exec.Command("polygon-edge",genCmd...)

	logWriter, err := os.Create(GenConfig.LogFile)
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