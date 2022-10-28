package inclusion

import (
	"math"

	"github.com/celestiaorg/celestia-app/pkg/shares"
	"github.com/celestiaorg/celestia-app/x/payment/types"
)

type path struct {
	instructions []WalkInstruction
	row          int
}

// calculateCommitPaths calculates all of the paths to subtree roots needed to
// create the commitment for a given message.
func calculateCommitPaths(squareSize, start, msgShareLen int) []path {
	// todo: make the non-interactive defaults optional. by calculating the
	// NextAlignedPowerOfTwo, we are forcing use of the non-interactive
	// defaults. If we want to make this optional in the future, we have to move
	// this next line out of this function.
	start, _ = shares.NextAlignedPowerOfTwo(start, msgShareLen, squareSize)
	msgMinSquareSize := MsgMinSquareSize(uint64(msgShareLen))
	startRow, endRow := start/squareSize, (start+msgShareLen-1)/squareSize
	normalizedStartIndex := start % squareSize
	// normalizedEndIndex is the index in the endRow where the message ends
	normalizedEndIndex := (start + msgShareLen) - endRow*squareSize
	paths := []path{}
	maxDepth := uint64(math.Log2(float64(squareSize)))
	for i := startRow; i <= endRow; i++ {
		start, end := 0, squareSize
		if i == startRow {
			start = normalizedStartIndex
		}
		if i == endRow {
			end = normalizedEndIndex
		}
		msgShareLenInRow := end - start
		coords := calculateSubTreeRootCoordinates(msgMinSquareSize, uint64(msgShareLenInRow), uint64(start))
		for _, coord := range coords {
			depth := maxDepth - coord.height
			paths = append(paths, path{
				instructions: genSubTreeRootPath(int(depth), uint(coord.position)),
				row:          i,
			})
		}
	}

	return paths
}

// genSubTreeRootPath calculates the path to a given subtree root of a node, given the
// depth and position of the node. note: the root of the tree is depth 0.
// The following nolint can be removed after this function is used.
//
//nolint:unused,deadcode
func genSubTreeRootPath(depth int, pos uint) []WalkInstruction {
	path := make([]WalkInstruction, depth)
	counter := 0
	for i := depth - 1; i >= 0; i-- {
		if (pos & (1 << i)) == 0 {
			path[counter] = WalkLeft
		} else {
			path[counter] = WalkRight
		}
		counter++
	}
	return path
}

// coord identifies a tree node using the height and position
//
//	Height       Position
//	3              0
//	              / \
//	             /   \
//	2           0     1
//	           /\     /\
//	1         0  1   2  3
//	         /\  /\ /\  /\
//	0       0 1 2 3 4 5 6 7
type coord struct {
	// height is the height of a node where the height of a leaf is 0
	height uint64
	// position is the index of a node at a given height, 0 being the left most
	// node
	position uint64
}

// calculateSubTreeRootCoordinates generates the sub tree root coordinates for
// the portion of a message in a row.
func calculateSubTreeRootCoordinates(msgMinSquareSize uint64, msgLengthInRow uint64, startIndex uint64) (result []coord) {
	index := startIndex
	treeSizes := types.MerkleMountainRangeSizes(msgLengthInRow, msgMinSquareSize)
	for _, treeSize := range treeSizes {
		height := uint64(math.Log2(float64(treeSize)))
		position := index / treeSize
		result = append(result, coord{
			height:   height,
			position: position,
		})
		index += treeSize
	}
	return result
}

// MsgMinSquareSize returns the minimum square size that msgLen can be included
// in. The returned square size does not account for the associated transaction
// shares or non-interactive defaults so it is a minimum.
func MsgMinSquareSize(msgLen uint64) uint64 {
	shareCount := shares.MsgSharesUsed(int(msgLen))
	squareSize := shares.RoundUpPowerOfTwo(int(math.Ceil(math.Sqrt(float64(shareCount)))))
	return uint64(squareSize)
}
