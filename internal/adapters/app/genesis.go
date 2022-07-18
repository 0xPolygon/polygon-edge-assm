package app

import (
	"fmt"
	"strings"
)

func (a *Adapter) generateGenesisCommand() []string {
	coreData := a.core.GetCore()

	// first part of genesis command
	genCmd := []string{"genesis", "--consensus", "ibft", "--dir", coreData.Config.FileLocation}

	// Name
	if coreData.Config.ChainName != "" {
		genCmd = append(genCmd, "--name", coreData.Config.ChainName)
	}

	// PoS params
	if coreData.Config.PoS {
		genCmd = append(genCmd, "--pos")
	}

	// Epoch Size
	if coreData.Config.EpochSize != "" {
		genCmd = append(genCmd, "--epoch-size", coreData.Config.EpochSize)
	}

	// Chain ID
	if coreData.Config.ChainID != "" {
		genCmd = append(genCmd, "--chain-id", coreData.Config.ChainID)
	}

	// Block Gas Limit
	if coreData.Config.BlockGasLimit != "" {
		genCmd = append(genCmd, "--block-gas-limit", coreData.Config.BlockGasLimit)
	}

	// Max validator count
	if coreData.Config.PoS && coreData.Config.MaxValidatorCount != "" {
		genCmd = append(genCmd, "--max-validator-count", coreData.Config.MaxValidatorCount)
	}

	// Min validator count
	if coreData.Config.PoS && coreData.Config.MinValidatorCount != "" {
		genCmd = append(genCmd, "--min-validator-count", coreData.Config.MinValidatorCount)
	}

	// add validators and keys
	for _, v := range coreData.AllNodesInitInfo {
		genCmd = append(genCmd, "--ibft-validator="+v.GenesisValidatorKey)
		genCmd = append(genCmd,
			fmt.Sprintf("--bootnode=/ip4/%s/tcp/1478/p2p/%s",
				strings.TrimSpace(v.IP),
				strings.TrimSpace(v.GenesisNetworkID)),
		)
	}

	// Premine
	for _, premine := range strings.Split(coreData.Config.Premine, ",") {
		genCmd = append(genCmd, fmt.Sprintf("--premine=%s", premine))
	}

	return genCmd
}
