package main

import (
	"testing"
)

func TestDeployDefaultNetwork(t *testing.T) {
	d := NewDeployer(DConfig())
	d.Deploy()
}
