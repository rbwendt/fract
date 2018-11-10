package state

import (
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"time"

	"github.com/fogleman/gg"
)

func doMandelbrot(state State, ctx *gg.Context, frame int) {

	fn := fmt.Sprintf("data/out-%04d.png", frame)

	if _, err := os.Stat(fn); !os.IsNotExist(err) {
		fmt.Println("file", fn, "exists, skipping")
		return
	}

	maxIteration := frame + 200
	if maxIteration < 500 {
		maxIteration = 500
	}

	prec := int(float64(frame)/10.0) + 21
	if prec < 63 {
		prec = 63
	}

	pts := getPts(state, maxIteration, prec)

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
		1920,
		1080}

	frames := 1805
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

		// fmt.Println(frame, &state.minX, &state.maxX, &state.minY, &state.maxY)

		f := func(frame int, state State) {
			if frame > 614 {
				// fmt.Println("computing frame", frame)
				ctx := gg.NewContext(state.w, state.h)
				start := time.Now()
				doMandelbrot(state, ctx, frame)
				elapsed := time.Since(start)
				fmt.Println("frame", frame, "took", elapsed)
			}
		}
		if frame%4 == 0 {
			f(frame, state)
		} else {
			go f(frame, state)
		}
		state = zoomAt(desiredCenterX, desiredCenterY, state, 0.95)
	}
}
