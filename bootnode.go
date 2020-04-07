package main

import (
	"crypto/ecdsa"
	"fmt"
	"net/url"
	"strconv"

	"github.com/ethereum/go-ethereum/crypto"
)

// GenSaveKey generates and saves a new key for bootnode.
func GenSaveKey(path string) (*ecdsa.PrivateKey, error) {
	nodeKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	return nodeKey, crypto.SaveECDSA(path, nodeKey)
}

// GetBootNodeString .
func GetBootNodeString(nodeKey *ecdsa.PrivateKey, host string, port int) string {
	u := url.URL{Scheme: "enode"}
	nodeid := fmt.Sprintf("%x", crypto.FromECDSAPub(&nodeKey.PublicKey)[1:])
	u.User = url.User(nodeid)
	u.Host = host + ":" + strconv.Itoa(port)
	return u.String()
}
