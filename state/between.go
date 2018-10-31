package state

import "math/big"

func between(a, b string, p float64) big.Float {
	prec := uint(206)

	base := 10
	c, _, _ := big.ParseFloat(a, base, prec, big.ToNearestEven)
	d, _, _ := big.ParseFloat(b, base, prec, big.ToNearestEven)

	q := big.NewFloat(1 - p)
	r := big.NewFloat(p)

	f := mul(c, q)
	g := mul(d, r)

	return add(&f, &g)
}
