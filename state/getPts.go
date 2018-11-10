package state
import (
	"math/big"
)

func getPts(state State, maxIteration, prec int) [][]int {

	pts := make([][]int, 0)

	xpts := state.w
	ypts := state.h
	minx := state.minX
	maxx := state.maxX
	miny := state.minY
	maxy := state.maxY
	xdif := sub(&maxx, &minx)
	ydif := sub(&maxy, &miny)

	xdif.SetPrec(uint(prec))
	ydif.SetPrec(uint(prec))

	var bxp big.Float
	bxp.SetInt(big.NewInt(int64(xpts)))

	var byp big.Float
	byp.SetInt(big.NewInt(int64(ypts)))

	for i := 0; i < xpts; i++ {

		var bi big.Float
		bi.SetInt(big.NewInt(int64(i)))
		row := make([]int, 0)
		for j := 0; j < ypts; j++ {

			var bj big.Float
			bj.SetInt(big.NewInt(int64(j)))
			xmul := mul(&xdif, &bi)
			ymul := mul(&ydif, &bj)

			xquo := quo(&xmul, &bxp)
			yquo := quo(&ymul, &byp)

			x0 := add(&xquo, &minx)
			y0 := add(&yquo, &miny)

			x := big.NewFloat(0.0)
			y := big.NewFloat(0.0)

			iteration := 0

			x2 := mul(x, x)
			y2 := mul(y, y)
			sum := add(&x2, &y2)

			for ; lt4(&sum) && (iteration < maxIteration); iteration++ {

				x2 = mul(x, x)
				y2 = mul(y, y)

				diff := sub(&x2, &y2)
				xtemp := add(&diff, &x0)

				twox := mul(big.NewFloat(2.0), x)
				twoxy := mul(&twox, y)

				ysum := add(&twoxy, &y0)

				y = &ysum
				x = &xtemp

				x2 = mul(x, x)
				y2 = mul(y, y)

				sum = add(&x2, &y2)
			}

			sum = add(&x2, &y2)
			if lt4(&sum) {
				iteration = 0
			}

			row = append(row, iteration)
		}
		pts = append(pts, row)

	}
	return pts
}