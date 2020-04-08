package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContact(t *testing.T) {
	cs := []struct {
		Input  []string
		Expect string
	}{
		{
			Input:  []string{"", ""},
			Expect: "",
		},
		{
			Input:  []string{"1", "2"},
			Expect: "1,2",
		},
		{
			Input:  []string{"1", ""},
			Expect: "1",
		},
		{
			Input:  []string{"", "1"},
			Expect: "1",
		},
	}

	for _, c := range cs {
		assert.Equal(t, c.Expect, Contact(c.Input))
	}
}
