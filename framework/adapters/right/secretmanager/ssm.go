package secretmanager

import (
	"context"
	"fmt"
	edgeCrypto "github.com/0xPolygon/polygon-edge/crypto"
	"github.com/0xPolygon/polygon-edge/helper/hex"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
)

type Adapter struct {
	region string
}

func NewAdapter(region string) *Adapter {
	return &Adapter{
		region: region,
	}
}

func (a Adapter) getSecret(secretName string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(a.region))
	if err != nil {
		return "", fmt.Errorf("could not load aws configurtation err=%w", err)
	}

	client := ssm.NewFromConfig(cfg)

	param, err := client.GetParameter(
		context.TODO(),
		&ssm.GetParameterInput{
			Name:           aws.String(secretName),
			WithDecryption: true})
	if err != nil {
		return "", fmt.Errorf("could not get the parameter from AWS SSM store err=%w", err)
	}

	return *param.Parameter.Value, nil
}

func (a Adapter) GetValidatorKey(key string) (string, error) {
	// get ssm stored key
	ssmValKey, err := a.getSecret(key)
	if err != nil {
		return "", fmt.Errorf("could not retreve validator key from ssm err=%w", err)
	}
	// Get the validator address from validator-key stored in AWS SSM
	valPrivKey, err := edgeCrypto.BytesToPrivateKey([]byte(ssmValKey))
	if err != nil {
		return "", fmt.Errorf("could not get validator address from private key: %w", err)
	}
	// get address from key
	valAddr, err := edgeCrypto.GetAddressFromKey(valPrivKey)
	if err != nil {
		return "", fmt.Errorf("could not get validator address from private key: %w", err)
	}

	return valAddr.String(), nil
}

func (a Adapter) GetNetworkKey(id string) (string, error) {
	ssmNetworkID, err := a.getSecret(id)
	if err != nil {
		return "", fmt.Errorf("could not retreve network id from ssm err=%w", err)
	}
	// get the libp2p network id from network private key
	buf, _ := hex.DecodeString(ssmNetworkID)

	networkPubKey, err := crypto.UnmarshalPrivateKey(buf)
	if err != nil {
		return "", fmt.Errorf("could not convert private to public network key: %w", err)
	}

	peerID, _ := peer.IDFromPrivateKey(networkPubKey)

	return peerID.String(), nil
}
