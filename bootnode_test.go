package main

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestBootNode(t *testing.T) {
	expect := "enode://7ce9736baa193db464b1f57a45320c66b645f8ec38140459af276f7064fce469ad562d75ac10746f9148d478fe750a6b8a4c1aeaa72b2032719c1713214297a5@localhost:30303"
	key, _ := crypto.HexToECDSA("7c4b9e50a61eba57d3ded2bcf4246f9e6cbc61166e2ce3525aea57b1a852f41b")
	assert.Equal(t, expect, GetBootNodeString(key, "localhost", 30303))
}
