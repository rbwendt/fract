package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/fogleman/gg"
)

type State struct {
	minX float64
	maxX float64
	minY float64
	maxY float64
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
	// fmt.Println(xpts, ypts, minx, miny, maxx, maxy)
	for i := 0; i < xpts; i++ {

		row := make([]int, 0)
		for j := 0; j < ypts; j++ {
			x0 := (maxx-minx)*float64(i)/float64(xpts) + minx
			y0 := (maxy-miny)*float64(j)/float64(ypts) + miny

			x := float64(0)
			xtemp := float64(0)
			y := float64(0)
			// fmt.Println(x0, y0, x, xtemp, y)
			iteration := 0

			for ; ((x*x)+(y*y) < 4) && (iteration < maxIteration); iteration++ {
				xtemp = x*x - y*y + x0
				y = 2*x*y + y0
				x = xtemp
			}

			if x*x+y*y < 4 {
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
	maxIteration := 1000

	pts := getPts(state, maxIteration)

	for i := 0; i < state.w; i++ {
		for j := 0; j < state.h; j++ {
			iteration := pts[i][j]
			var l float64
			l = 0.5
			// if iteration < 20 {
			// 	l = float64(iteration)
			// } else {
			// 	l = 50
			// }

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

func zoomAt(centerX, centerY float64, state State, zoomFactor float64) State {
	return State{

		centerX - zoomFactor*(state.maxX-state.minX)/2,
		centerX + zoomFactor*(state.maxX-state.minX)/2,

		centerY - zoomFactor*(state.maxY-state.minY)/2,
		centerY + zoomFactor*(state.maxY-state.minY)/2,
		state.w,
		state.h}
}

func main() {
	exec.Command("rm data/*.png")
	state := State{-2.5, 1, -1, 1, 640, 480}
	frames := 1000
	desiredCenterX := -1.235883546447
	desiredCenterY := -0.1091632575204
	fmt.Println(desiredCenterX, desiredCenterY)
	for frame := 0; frame < frames; frame++ {
		fmt.Println(frame, state)
		ctx := gg.NewContext(state.w, state.h)
		start := time.Now()
		if frame > 570 {
			doMandelbrot(state, ctx, frame)
		}
		elapsed := time.Since(start)
		fmt.Println("frame", frame, "took", elapsed)
		state = zoomAt(desiredCenterX, desiredCenterY, state, 0.95)
	}
}
