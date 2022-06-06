package common

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/pkg/errors"
)

// Node represents a node in a Tree.
type Node struct {
	key       []byte
	value     []byte
	Hash      []byte
	LeftHash  []byte
	RightHash []byte
	Version   int64
	size      int64
	height    int8
}

// NewNode returns a new node from a key, value and version.
func NewNode(key []byte, value []byte, version int64) *Node {
	return &Node{
		key:     key,
		value:   value,
		height:  0,
		size:    1,
		Version: version,
	}
}

// MakeNode constructs an *Node from an encoded byte slice.
//
// The new node doesn't have its hash saved or set. The caller must set it
// afterwards.
func MakeNode(buf []byte) (*Node, error) {

	// Read node header (height, size, version, key).
	height, n, cause := decodeVarint(buf)
	if cause != nil {
		return nil, errors.Wrap(cause, "decoding node.height")
	}
	buf = buf[n:]
	if height < int64(math.MinInt8) || height > int64(math.MaxInt8) {
		return nil, errors.New("invalid height, must be int8")
	}

	size, n, cause := decodeVarint(buf)
	if cause != nil {
		return nil, errors.Wrap(cause, "decoding node.size")
	}
	buf = buf[n:]

	ver, n, cause := decodeVarint(buf)
	if cause != nil {
		return nil, errors.Wrap(cause, "decoding node.version")
	}
	buf = buf[n:]

	key, n, cause := decodeBytes(buf)
	if cause != nil {
		return nil, errors.Wrap(cause, "decoding node.key")
	}
	buf = buf[n:]

	node := &Node{
		height:  int8(height),
		size:    size,
		Version: ver,
		key:     key,
	}

	// Read node body.

	if node.isLeaf() {
		val, _, cause := decodeBytes(buf)
		if cause != nil {
			return nil, errors.Wrap(cause, "decoding node.value")
		}
		node.value = val
	} else { // Read children.
		leftHash, n, cause := decodeBytes(buf)
		if cause != nil {
			return nil, errors.Wrap(cause, "deocding node.leftHash")
		}
		buf = buf[n:]

		rightHash, _, cause := decodeBytes(buf)
		if cause != nil {
			return nil, errors.Wrap(cause, "decoding node.rightHash")
		}
		node.LeftHash = leftHash
		node.RightHash = rightHash
	}
	return node, nil
}

func (node *Node) isLeaf() bool {
	return node.height == 0
}

// decodeBytes decodes a varint length-prefixed byte slice, returning it along with the number
// of input bytes read.
func decodeBytes(bz []byte) ([]byte, int, error) {
	s, n, err := decodeUvarint(bz)
	if err != nil {
		return nil, n, err
	}
	// Make sure size doesn't overflow. ^uint(0) >> 1 will help determine the
	// max int value variably on 32-bit and 64-bit machines. We also doublecheck
	// that size is positive.
	size := int(s)
	if s >= uint64(^uint(0)>>1) || size < 0 {
		return nil, n, fmt.Errorf("invalid out of range length %v decoding []byte", s)
	}
	// Make sure end index doesn't overflow. We know n>0 from decodeUvarint().
	end := n + size
	if end < n {
		return nil, n, fmt.Errorf("invalid out of range length %v decoding []byte", size)
	}
	// Make sure the end index is within bounds.
	if len(bz) < end {
		return nil, n, fmt.Errorf("insufficient bytes decoding []byte of length %v", size)
	}
	bz2 := make([]byte, size)
	copy(bz2, bz[n:end])
	return bz2, end, nil
}

// decodeUvarint decodes a varint-encoded unsigned integer from a byte slice, returning it and the
// number of bytes decoded.
func decodeUvarint(bz []byte) (uint64, int, error) {
	u, n := binary.Uvarint(bz)
	if n == 0 {
		// buf too small
		return u, n, errors.New("buffer too small")
	} else if n < 0 {
		// value larger than 64 bits (overflow)
		// and -n is the number of bytes read
		n = -n
		return u, n, errors.New("EOF decoding uvarint")
	}
	return u, n, nil
}

// decodeVarint decodes a varint-encoded integer from a byte slice, returning it and the number of
// bytes decoded.
func decodeVarint(bz []byte) (int64, int, error) {
	i, n := binary.Varint(bz)
	if n == 0 {
		return i, n, errors.New("buffer too small")
	} else if n < 0 {
		// value larger than 64 bits (overflow)
		// and -n is the number of bytes read
		n = -n
		return i, n, errors.New("EOF decoding varint")
	}
	return i, n, nil
}
