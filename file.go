package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/core"
)

// FlushGenesis dumps genesis json config to disk.
func FlushGenesis(genesis *core.Genesis, path string) {
	os.MkdirAll(filepath.Dir(path), 0755)

	out, _ := json.MarshalIndent(genesis, "", "  ")
	if err := Flush(out, path); err != nil {
		panic(fmt.Sprintln("Failed to save puppeth configs", "file", path, "err", err))
	}
}

// Flush dumps contents to disk.
func Flush(out []byte, path string) error {
	return ioutil.WriteFile(path, out, 0644)
}
