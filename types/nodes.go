package types

import (
	"encoding/hex"
	"fmt"

	edgeCrypto "github.com/0xPolygon/polygon-edge/crypto"
	"github.com/libp2p/go-libp2p-core/crypto"
	peer "github.com/libp2p/go-libp2p-core/peer"
)

type Responce struct {
	Success bool `json:"success"`
	Message string `json:"err_msg"`
}

type NodeInfo struct {
	NetworkID string
	ValidatorKey string
	IP string
}

type Nodes struct {
	Total int
	Finished []string
	NodeIPs map[string]string
	Node map[string]NodeInfo
}


func (n *NodeInfo) initValidatorKey(key string) error {
	// Get the validator address from validator-key stored in AWS SSM
	valPrivKey, err := edgeCrypto.BytesToPrivateKey([]byte(key))
	if err != nil {
		return fmt.Errorf("could not get validator address from private key: %w", err)
	}

	valAddr, err := edgeCrypto.GetAddressFromKey(valPrivKey)
	if err != nil {
		return fmt.Errorf("could not get validator address from private key: %w", err)
	}
	n.ValidatorKey = valAddr.String()
 
	return nil
}

func (n *NodeInfo) initNetworkKey(id string) error {
	// get the libp2p network id from network private key
	buf, _ := hex.DecodeString(id)
	networkPubKey, err := crypto.UnmarshalPrivateKey(buf)
	if err != nil {
		return fmt.Errorf("could not convert private to public network key: %w", err)
	}
	
	peerId, _ := peer.IDFromPrivateKey(networkPubKey)
	n.NetworkID = peerId.String() 
	
	return nil
}

func NewNodeInfo(networkIdPrivKey string, validatorPrivKey string, ipAddress string) (*NodeInfo, error) {
	nInfo := &NodeInfo{}
	
	if err := nInfo.initNetworkKey(networkIdPrivKey); err != nil {
		return nil, fmt.Errorf("could not set network id: %w", err)
	}

	if err := nInfo.initValidatorKey(validatorPrivKey); err != nil {
		return nil, fmt.Errorf("could not set validator key: %w", err)
	}

	nInfo.IP = ipAddress
	
	return nInfo, nil
}


