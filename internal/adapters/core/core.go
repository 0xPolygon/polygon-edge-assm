package core

import (
	"encoding/json"
)

type Core struct {
	Nodes  `json:"nodes,omitempty"`
	Config Config `json:"config,omitempty"`
}

// Nodes is the struct holding all the information about nodes
type Nodes struct {
	Total              int            `json:"total"`
	AllNodesInitInfo   []NodeInitInfo `json:"all_node_data"`
	SingleNodeInitInfo NodeInitInfo   `json:"node_info"`
}

// NodeInitInfo is the struct holding the network information about nodes
type NodeInitInfo struct {
	SSMParamID string `json:"ssm_param_id"`
	IP         string `json:"ip"`
	NodeName   string `json:"node_name"`

	GenesisValidatorKey string
	GenesisNetworkID    string
}

// Config holds the general info
type Config struct {
	AWSRegion    string `json:"aws_region"`
	S3BucketName string `json:"s3_bucket_name"`
	S3KeyName    string `json:"s3_key_name"`

	GenesisConfig
}

// GenesisConfig holds the info about genesis
type GenesisConfig struct {
	ChainName         string `json:"chain_name"`
	ChainID           string `json:"chain_id"`
	Premine           string `json:"premine"`
	PoS               bool   `json:"pos"`
	EpochSize         string `json:"epoch_size"`
	BlockGasLimit     string `json:"block_gas_limit"`
	MaxValidatorCount string `json:"max_validator_count"`
	MinValidatorCount string `json:"min_validator_count"`
	Consensus         string `json:"consensus"`

	FileLocation string
}

func (c *Config) ToJSON() string {
	jsonData, _ := json.Marshal(c)

	return string(jsonData)
}

func NewAdapter() *Core {
	return &Core{
		Nodes:  Nodes{},
		Config: Config{},
	}
}

func (a *Core) GetConfig() *Config {
	return &a.Config
}

func (a *Core) GetNodesInfo() *Nodes {
	return &a.Nodes
}

func (a *Core) GetCore() *Core {
	return a
}

func (a *Core) GetCoreJSON() string {
	jsonData, _ := json.Marshal(a)

	return string(jsonData)
}
