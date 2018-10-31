package state

import "math/big"

func quo(a, b *big.Float) big.Float {
	var d big.Float
	d = *(d.Quo(a, b))
	return d
}

func mul(a, b *big.Float) big.Float {
	var d big.Float
	d = *(d.Mul(a, b))
	return d
}

func add(a, b *big.Float) big.Float {
	var d big.Float
	d = *(d.Add(a, b))
	return d
}

func sub(a, b *big.Float) big.Float {
	var d big.Float
	d = *(d.Sub(a, b))
	return d
}

func lt4(a *big.Float) bool {
	four := big.NewFloat(4.0)
	return lt(a, four)
}

func lt(a, b *big.Float) bool {
	return (a.Cmp(b) == -1)
}
