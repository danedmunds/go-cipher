package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Buffer struct {
	bytes.Buffer
}

func (b *Buffer) Close() error {
	return nil
}

func TestRailFenceEncipher(t *testing.T) {
	num := 3
	rails := make([]io.ReadWriteCloser, 0, num)
	for i := 0; i < num; i++ {
		rails = append(rails, &Buffer{})
	}
	railFence := NewRailFencerEncipher(rails)

	for _, r := range []rune("WEAREDISCOVEREDFLEEATONCE") {
		railFence.Write(r)
	}

	out := &bytes.Buffer{}
	railFence.DumpTo(out)

	assert.Equal(t, "WECRLTEERDSOEEFEAOCAIVDEN", string(out.Bytes()))
}

func TestSelectRail(t *testing.T) {
	tests := []struct {
		numRails int
		sequence []int
	}{
		{
			1,
			[]int{0, 0, 0, 0},
		},
		{
			2,
			[]int{0, 1, 0, 1},
		},
		{
			3,
			[]int{0, 1, 2, 1, 0},
		},
		{
			4,
			[]int{0, 1, 2, 3, 2, 1, 0},
		},
		{
			5,
			[]int{0, 1, 2, 3, 4, 3, 2, 1, 0},
		},
	}

	for _, test := range tests {
		for i, rail := range test.sequence {
			assert.Equal(t, selectRail(test.numRails, i), rail, "num rails: %d, index: %d", test.numRails, i)
		}
	}
}
