package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/core"
)

// flushGenesis dumps genesis json config to disk.
func flushGenesis(genesis *core.Genesis, path string) {
	os.MkdirAll(filepath.Dir(path), 0755)

	out, _ := json.MarshalIndent(genesis, "", "  ")
	if err := Flush(out, path); err != nil {
		panic(fmt.Sprintln("Failed to save puppeth configs", "file", path, "err", err))
	}
}

// StoreKey stores key in keystore format.
func StoreKey(key *ecdsa.PrivateKey, dir, path string) {
	nks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)
	nks.ImportECDSA(key, pw)
}

// Flush dumps contents to disk.
func Flush(out []byte, path string) error {
	return ioutil.WriteFile(path, out, 0644)
}
