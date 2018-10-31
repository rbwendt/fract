package state

import (
	"fmt"
	"math/big"
	"os/exec"
	"time"

	"github.com/fogleman/gg"
)

func getPts(state State, maxIteration int) [][]int {

	pts := make([][]int, 0)

	xpts := state.w
	ypts := state.h
	minx := state.minX
	maxx := state.maxX
	miny := state.minY
	maxy := state.maxY
	xdif := sub(&maxx, &minx)
	ydif := sub(&maxy, &miny)

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

func doMandelbrot(state State, ctx *gg.Context, frame int) {
	maxIteration := 2000

	pts := getPts(state, maxIteration)

	for i := 0; i < state.w; i++ {
		for j := 0; j < state.h; j++ {
			iteration := pts[i][j]
			l := 0.5

			h := float64(iteration) / float64(maxIteration)
			// h = circleify(h)
			if iteration < 20 {
				l = float64(iteration) / 100.0
			}
			r, g, b := hsl2rgb(h, 1, l)
			ctx.SetRGB(r, g, b)

			ctx.SetPixel(i, j)
			ctx.Fill()
		}
	}

	fn := fmt.Sprintf("data/out-%04d.png", frame)
	ctx.SavePNG(fn)
}

// this makes it go.
func Go() {
	exec.Command("rm data/*.png")
	state := State{
		*big.NewFloat(-2.5),
		*big.NewFloat(1),
		*big.NewFloat(-1),
		*big.NewFloat(1),
		640,
		480}

	frames := 10000
	// yc := "-0.10935"

	// desiredCenterX, _, _ := big.ParseFloat(xc, base, prec, big.ToNearestEven)
	desiredCenterX := between(
		"-1.2369998000011500984940384605701753937082146811457768156",
		"-1.23699980000115009849403846057015875884176836656296463",
		0.8)
	// desiredCenterY, _, _ := big.ParseFloat(yc, base, prec, big.ToNearestEven)
	desiredCenterY := between(
		"-0.1093500000000000000000000000028940139019990971",
		"-0.1093499999999999999999999999971059860980009029",
		0.8)
	// desiredCenterY := big.NewFloat(-0.109350000010000)
	// I have to set up parsing for these longer inputs.

	for frame := 0; frame < frames; frame++ {
		ctx := gg.NewContext(state.w, state.h)
		if frame > 1800 {
			fmt.Println(frame, &state.minX, &state.maxX, &state.minY, &state.maxY)

			start := time.Now()
			doMandelbrot(state, ctx, frame)
			elapsed := time.Since(start)
			fmt.Println("frame took", elapsed)
			break
		}

		state = zoomAt(desiredCenterX, desiredCenterY, state, 0.95)
	}
}
