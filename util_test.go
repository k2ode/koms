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

func TestUtilMakeInc(t *testing.T) {
	inc := MakeInc(5)

	assert.Equal(t, inc(1), 2)
	assert.Equal(t, inc(4), 5)
	assert.Equal(t, inc(5), 0)
}

func TestUtilMakeDesc(t *testing.T) {
	desc := MakeDesc(5)

	assert.Equal(t, desc(5), 4)
	assert.Equal(t, desc(1), 0)
	assert.Equal(t, desc(0), 5)
}