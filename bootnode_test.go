package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBootNode(t *testing.T) {
	key, err := GenSaveKey("boot.key")
	assert.Nil(t, err)

	t.Log(GetBootNodeString(key, "localhost", 30303))
}