package main

import "testing"

func TestConfig(t *testing.T) {
	WriteConfig("test.yaml", DConfig())
}
