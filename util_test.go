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

func TestUtilMakeIncBy(t *testing.T) {
	inc := MakeIncBy(5, 2)

	assert.Equal(t, inc(0), 2)
	assert.Equal(t, inc(3), 5)

	assert.Equal(t, inc(5), 0)

	assert.Equal(t, inc(3), 5)
	assert.Equal(t, inc(4), 5)
}


func TestUtilMakeDescBy(t *testing.T) {
	desc := MakeDescBy(5, 2)

	assert.Equal(t, desc(5), 3)
	assert.Equal(t, desc(4), 2)

	assert.Equal(t, desc(0), 5)

	assert.Equal(t, desc(2), 0)
	assert.Equal(t, desc(1), 0)
}

// func TestUtilFind(t *testing.T) {
// 	arr := []int{1, 2, 3, 4, 5}

// 	assert.Equal(t, Find(arr, 0), -1)
// 	assert.Equal(t, Find(arr, 1), 0)
// 	assert.Equal(t, Find(arr, 5), 4)
// }

// func TestUtilContains(t *testing.T) {
// 	arr := []int{1, 2, 3, 4, 5}

// 	assert.False(t, Contains(arr, 0))
// 	assert.True(t, Contains(arr, 1))
// 	assert.True(t, Contains(arr, 5))
// }