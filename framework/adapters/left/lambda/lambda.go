package lambda

import (
	"Trapesys/polygon-edge-assm/internal/adapters/core"
)

type Adapter struct {
	config core.Config
	nodes  core.Nodes
}

func NewAdapter() *Adapter {
	return &Adapter{
		config: core.Config{},
		nodes:  core.Nodes{},
	}
}

func (a *Adapter) SetNodes(config core.Core) {
	a.nodes = core.Nodes{
		Total: config.Total,
		SingleNodeInitInfo: core.NodeInitInfo{
			SSMParamID: config.SingleNodeInitInfo.SSMParamID,
			IP:         config.SingleNodeInitInfo.IP,
			NodeName:   config.SingleNodeInitInfo.NodeName,
		},
	}
}

func (a Adapter) GetNodes() core.Nodes {
	return a.nodes
}

func (a *Adapter) SetConfig(config core.Core) {
	a.config.ChainName = config.Config.ChainName
	a.config.Premine = config.Config.Premine
	a.config.Consensus = config.Config.Consensus
	a.config.BlockGasLimit = config.Config.BlockGasLimit
	a.config.ChainID = config.Config.ChainID
	a.config.EpochSize = config.Config.EpochSize
	a.config.MaxValidatorCount = config.Config.MaxValidatorCount
	a.config.MinValidatorCount = config.Config.MinValidatorCount
	a.config.PoS = config.Config.PoS

	a.config.AWSRegion = config.Config.AWSRegion
	a.config.S3BucketName = config.Config.S3BucketName
	a.config.S3KeyName = config.Config.S3KeyName
}

func (a Adapter) GetConfig() core.Config {
	return a.config
}
