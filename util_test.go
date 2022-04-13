package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtilLines(t *testing.T) {
	str := "foo\nbar"
	lines := Lines(str)
	assert.Equal(t, lines[0], "foo")
	assert.Equal(t, lines[1], "bar")
}

func TestUtilBetween(t *testing.T) {
	between := Between("what is a foo bar", "what is a ", " bar")

	assert.Equal(t, between, "foo")

}