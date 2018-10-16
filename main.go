package main

import (
	"fmt"
	"math/big"
	"os/exec"

	"github.com/fogleman/gg"
)

type State struct {
	minX big.Float
	maxX big.Float
	minY big.Float
	maxY big.Float
	w    int
	h    int
}

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

// borrowed https://gist.github.com/emanuel-sanabria-developer/5793377
func hue2rgb(p, q, t float64) float64 {
	if t < 0 {
		t++
	}
	if t > 1 {
		t--
	}
	if t < 0.1666666666 {
		return p + (q-p)*6.0*t
	}
	if t < 0.5 {
		return q
	}
	if t < 0.666666666 {
		return p + (q-p)*(0.666666666-t)*6.0
	}
	return p
}

func hsl2rgb(h, s, l float64) (r, g, b float64) {
	if s == 0 {
		r = l
		g = l
		b = l // achromatic
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}

		p := 2.0*l - q
		r = hue2rgb(p, q, h+0.33333)
		g = hue2rgb(p, q, h)
		b = hue2rgb(p, q, h-0.3333333)
	}

	return r, g, b
}

func doMandelbrot(state State, ctx *gg.Context, frame int) {
	maxIteration := 500

	pts := getPts(state, maxIteration)

	for i := 0; i < state.w; i++ {
		for j := 0; j < state.h; j++ {
			iteration := pts[i][j]
			l := 0.5

			h := float64(iteration) / float64(maxIteration)
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

func scale(max, min big.Float, zoomFactor float64) *big.Float {
	z := big.NewFloat(0.5 * zoomFactor)
	diff := sub(&max, &min)
	res := mul(z, &diff)
	return &res
}

func zoomAt(centerX, centerY big.Float, state State, zoomFactor float64) State {

	xsub := sub(&centerX, scale(state.maxX, state.minX, zoomFactor))
	xadd := add(&centerX, scale(state.maxX, state.minX, zoomFactor))

	ysub := sub(&centerY, scale(state.maxY, state.minY, zoomFactor))
	yadd := add(&centerY, scale(state.maxY, state.minY, zoomFactor))

	return State{

		xsub,
		xadd,

		ysub,
		yadd,

		state.w,
		state.h}
}

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

func main() {
	exec.Command("rm data/*.png")
	state := State{
		*big.NewFloat(-2.5),
		*big.NewFloat(1),
		*big.NewFloat(-1),
		*big.NewFloat(1),
		640,
		480}

	frames := 10
	desiredCenterX := big.NewFloat(-1.235883546447)
	desiredCenterY := big.NewFloat(-0.1091632575204)
	for frame := 0; frame < frames; frame++ {
		fmt.Println(frame, state)
		ctx := gg.NewContext(state.w, state.h)

		doMandelbrot(state, ctx, frame)

		state = zoomAt(*desiredCenterX, *desiredCenterY, state, 0.95)
	}
}
