package state

import "math/big"

type State struct {
	minX big.Float
	maxX big.Float
	minY big.Float
	maxY big.Float
	w    int
	h    int
}
