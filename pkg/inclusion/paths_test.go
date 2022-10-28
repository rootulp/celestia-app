package inclusion

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_calculateSubTreeRootCoordinates(t *testing.T) {
	type test struct {
		name             string
		msgMinSquareSize uint64
		msgLengthInRow   uint64
		startIndex       uint64
		expected         []coord
	}
	tests := []test{
		{
			name:             "first 4 shares of an 8 leaf tree",
			msgMinSquareSize: 2,
			msgLengthInRow:   4,
			startIndex:       0,
			expected: []coord{
				{
					height:   1,
					position: 0,
				},
				{
					height:   1,
					position: 1,
				},
			},
		},
		{
			name:             "last 4 shares of an 8 leaf tree",
			msgMinSquareSize: 2,
			msgLengthInRow:   4,
			startIndex:       4,
			expected: []coord{
				{
					height:   1,
					position: 2,
				},
				{
					height:   1,
					position: 3,
				},
			},
		},
		{
			name:             "last 4 shares of an 8 leaf tree",
			msgMinSquareSize: 2,
			msgLengthInRow:   4,
			startIndex:       4,
			expected: []coord{
				{
					height:   1,
					position: 2,
				},
				{
					height:   1,
					position: 3,
				},
			},
		},
	}
	for _, tt := range tests {
		res := calculateSubTreeRootCoordinates(tt.msgMinSquareSize, tt.msgLengthInRow, tt.startIndex)
		assert.Equal(t, tt.expected, res, tt.name)
	}
}

func Test_genSubTreeRootPath(t *testing.T) {
	type test struct {
		depth    int
		pos      uint
		expected []WalkInstruction
	}
	tests := []test{
		{2, 0, []WalkInstruction{WalkLeft, WalkLeft}},
		{0, 0, []WalkInstruction{}},
		{3, 0, []WalkInstruction{WalkLeft, WalkLeft, WalkLeft}},
		{3, 1, []WalkInstruction{WalkLeft, WalkLeft, WalkRight}},
		{3, 2, []WalkInstruction{WalkLeft, WalkRight, WalkLeft}},
		{5, 16, []WalkInstruction{WalkRight, WalkLeft, WalkLeft, WalkLeft, WalkLeft}},
	}
	for _, tt := range tests {
		path := genSubTreeRootPath(tt.depth, tt.pos)
		assert.Equal(t, tt.expected, path)
	}
}

func Test_calculateCommitPaths(t *testing.T) {
	type test struct {
		size, start, msgLen int
		expected            []path
	}
	tests := []test{
		{2, 0, 1, []path{{instructions: []WalkInstruction{WalkLeft}, row: 0}}},
		{2, 2, 2, []path{{instructions: []WalkInstruction{}, row: 1}}},
		{2, 1, 2, []path{{instructions: []WalkInstruction{}, row: 1}}},
		{4, 2, 2, []path{{instructions: []WalkInstruction{WalkRight}, row: 0}}},
		{4, 2, 4, []path{{instructions: []WalkInstruction{}, row: 1}}},
		{4, 3, 4, []path{{instructions: []WalkInstruction{}, row: 1}}},
		{4, 2, 9, []path{
			{instructions: []WalkInstruction{}, row: 1},
			{instructions: []WalkInstruction{}, row: 2},
			{instructions: []WalkInstruction{WalkLeft, WalkLeft}, row: 3},
		}},
		{8, 3, 16, []path{
			{instructions: []WalkInstruction{}, row: 1},
			{instructions: []WalkInstruction{}, row: 2},
		}},
		{64, 144, 32, []path{
			{instructions: []WalkInstruction{WalkRight}, row: 2},
		}},
		{64, 4032, 33, []path{
			{instructions: []WalkInstruction{WalkLeft}, row: 63},
			{instructions: []WalkInstruction{WalkRight, WalkLeft, WalkLeft, WalkLeft, WalkLeft, WalkLeft}, row: 63},
		}},
		{64, 4032, 63, []path{
			{instructions: []WalkInstruction{WalkLeft}, row: 63},
			{instructions: []WalkInstruction{WalkRight, WalkLeft}, row: 63},
			{instructions: []WalkInstruction{WalkRight, WalkRight, WalkLeft}, row: 63},
			{instructions: []WalkInstruction{WalkRight, WalkRight, WalkRight, WalkLeft}, row: 63},
			{instructions: []WalkInstruction{WalkRight, WalkRight, WalkRight, WalkRight, WalkLeft}, row: 63},
			{instructions: []WalkInstruction{WalkRight, WalkRight, WalkRight, WalkRight, WalkRight, WalkLeft}, row: 63},
		}},
	}
	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("test %d: square size %d start %d msgLen %d", i, tt.size, tt.start, tt.msgLen),
			func(t *testing.T) {
				assert.Equal(t, tt.expected, calculateCommitPaths(tt.size, tt.start, tt.msgLen))
			},
		)
	}
}
