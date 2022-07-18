package app

import (
	"Trapesys/polygon-edge-assm/framework/adapters/right/s3storage"
	"Trapesys/polygon-edge-assm/framework/adapters/right/secretmanager"
	"Trapesys/polygon-edge-assm/framework/ports"
	"Trapesys/polygon-edge-assm/internal/adapters/core"
	ports2 "Trapesys/polygon-edge-assm/internal/ports"
	"encoding/json"
	"fmt"
	"os"
)

type Adapter struct {
	lambdaAPI    ports.ILambdaAPIPort
	core         ports2.ICore
	s3Storage    ports.S3Storage
	assm         ports.IAwsSSMPort
	localStorage ports.ILocalStoragePort
}

func (a *Adapter) LambdaHandler(request core.Core) (string, error) {
	// get data from lambda api
	a.lambdaAPI.SetConfig(request)
	a.lambdaAPI.SetNodes(request)
	lambdaConfigData := a.lambdaAPI.GetConfig()
	lambdaNodesData := a.lambdaAPI.GetNodes()

	// set config data to internal structure
	if err := a.SetConfig(lambdaConfigData); err != nil {
		return "", fmt.Errorf("could not set up config err=%w", err)
	}

	// get node config and config from internal structure
	conf := a.core.GetConfig()
	// storage needs to be instantiated here as we need to set up a region
	s3Api, err := s3storage.NewAdapter(conf.AWSRegion, conf.S3BucketName)
	if err != nil {
		return "", fmt.Errorf("could not create new S3 adapter err=%w", err)
	}

	a.s3Storage = s3Api
	// storage needs to be instantiated here as we need to set up a region
	a.assm = secretmanager.NewAdapter(conf.AWSRegion)

	// fetch existing node data
	existingNodeData, err := a.s3Storage.FetchData(conf.S3KeyName)
	if err != nil {
		return "", fmt.Errorf("could not fetch data from S3 bucket err=%w", err)
	}

	// set data regarding the node info internal structure based on already existing data
	if err := a.SetNodes(lambdaNodesData, existingNodeData); err != nil {
		return "", fmt.Errorf("could not set node data internal structure: %w", err)
	}

	nodesInfo := a.core.GetNodesInfo()

	// If total nodes is equal to the len of all nodes we can generate genesis.json
	if nodesInfo.Total == len(nodesInfo.AllNodesInitInfo) {
		// fetch polygon-edge binary
		if err := a.localStorage.GetEdge(); err != nil {
			return "", fmt.Errorf("could not set local polygon-edge binary err=%w", err)
		}

		// fetch ssm stored keys and translate them to genesis accepted format
		for i, node := range nodesInfo.AllNodesInitInfo {
			valKey, err := a.assm.GetValidatorKey(fmt.Sprintf("/%s/%s/validator-key", node.SSMParamID, node.NodeName))
			if err != nil {
				return "", err
			}

			netwKey, err := a.assm.GetNetworkKey(fmt.Sprintf("/%s/%s/network-key", node.SSMParamID, node.NodeName))
			if err != nil {
				return "", err
			}

			// set validator and network keys to the internal structure
			newNode := core.NodeInitInfo{
				SSMParamID:          node.SSMParamID,
				IP:                  node.IP,
				NodeName:            node.NodeName,
				GenesisValidatorKey: valKey,
				GenesisNetworkID:    netwKey,
			}

			nodesInfo.AllNodesInitInfo[i] = newNode
		}

		// generate genesis command and create genesis file
		if err := a.localStorage.RunGenesisCmd(a.generateGenesisCommand()); err != nil {
			return "", fmt.Errorf("could not create genesis file err=%w", err)
		}

		// upload genesis.json to S3
		genesisFile, err := os.ReadFile(a.core.GetCore().Config.FileLocation)
		if err != nil {
			return "", fmt.Errorf("could not read genesis.json file err=%w", err)
		}

		if err := a.s3Storage.WriteData("genesis.json", string(genesisFile)); err != nil {
			return "", fmt.Errorf("could not write genesis.json to S3 err=%w", err)
		}

		return "Genesis file successfully created and uploaded to S3", nil
	}

	// write data to s3
	s3WriteErr := a.s3Storage.WriteData(conf.S3KeyName, a.core.GetCoreJSON())
	if s3WriteErr != nil {
		return "", fmt.Errorf("could not write data to S3 err=%w", err)
	}

	return "Node information successfully saved", nil
}

func (a *Adapter) SetNodes(recivedConf core.Nodes, existingNodeData string) error {
	// get the core structure pointer
	coreStruct := a.core.GetCore()
	// set total node info
	coreStruct.Total = recivedConf.Total

	// set all node info from json input saved in s3 from previous function run
	if err := json.Unmarshal([]byte(existingNodeData), coreStruct); err != nil {
		return fmt.Errorf("could not unmarshal exixting nodes json data err=%w", err)
	}

	// append this new request data to the existing all node info
	coreStruct.AllNodesInitInfo = append(coreStruct.AllNodesInitInfo, core.NodeInitInfo{
		SSMParamID: recivedConf.SingleNodeInitInfo.SSMParamID,
		IP:         recivedConf.SingleNodeInitInfo.IP,
		NodeName:   recivedConf.SingleNodeInitInfo.NodeName,
	})

	return nil
}

func (a *Adapter) SetConfig(receivedConf core.Config) error {
	if receivedConf.AWSRegion == "" {
		return fmt.Errorf("aws_region parameter missing")
	}

	if receivedConf.Premine == "" {
		return fmt.Errorf("premine address missing")
	}

	if receivedConf.S3BucketName == "" {
		return fmt.Errorf("s3 bucket name parameter missing")
	}

	if receivedConf.S3KeyName == "" {
		receivedConf.S3KeyName = "polygon-edge-config"
	}

	conf := a.core.GetConfig()
	// general config
	conf.AWSRegion = receivedConf.AWSRegion
	conf.S3KeyName = receivedConf.S3KeyName
	conf.S3BucketName = receivedConf.S3BucketName
	// genesis config
	conf.ChainName = receivedConf.ChainName
	conf.Premine = receivedConf.Premine
	conf.PoS = receivedConf.PoS
	conf.Consensus = receivedConf.Consensus
	conf.MinValidatorCount = receivedConf.MinValidatorCount
	conf.MaxValidatorCount = receivedConf.MaxValidatorCount
	conf.EpochSize = receivedConf.EpochSize
	conf.ChainID = receivedConf.ChainID
	conf.BlockGasLimit = receivedConf.BlockGasLimit

	conf.FileLocation = "/tmp/genesis.json"

	return nil
}

// NewAdapter injects lambdaAPI adapter that is used to send responses
func NewAdapter(core *core.Core, lambdaRight ports.ILambdaAPIPort, localStorage ports.ILocalStoragePort) *Adapter {
	return &Adapter{
		lambdaAPI:    lambdaRight,
		core:         core,
		localStorage: localStorage,
	}
}
