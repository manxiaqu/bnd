package main

import (
	"testing"
)

func TestMakeGenesis(t *testing.T) {
	ethashPath := "testdata/ethash.json"
	MakeGenesis(testEthashGenesisConfig(), ethashPath)

	cliquePath := "testdata/clique.json"
	MakeGenesis(testCliqueGenesisConfig(), cliquePath)
}
