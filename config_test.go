package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigAssign(t *testing.T) {
	c := &Config{}
	c.AssignDefault()

	assert.Equal(t, c, DConfig())
}
