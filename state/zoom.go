package state

import "math/big"

func scale(max, min big.Float, zoomFactor float64) big.Float {
	z := big.NewFloat(0.5 * zoomFactor)
	diff := sub(&max, &min)
	res := mul(z, &diff)
	// fmt.Println("z", z, "max", &max, "min", &min, "diff", &diff, "res", &res)
	return res
}

func zoomAt(centerX, centerY big.Float, state State, zoomFactor float64) State {

	xscale := scale(state.maxX, state.minX, zoomFactor)
	yscale := scale(state.maxY, state.minY, zoomFactor)

	xsub := sub(&centerX, &xscale)
	xadd := add(&centerX, &xscale)

	ysub := sub(&centerY, &yscale)
	yadd := add(&centerY, &yscale)

	// fmt.Println("xscale", &xscale, "xsub", &xsub, "xadd", &xadd, "yscale", &yscale)

	return State{
		xsub,
		xadd,

		ysub,
		yadd,

		state.w,
		state.h}
}
