package state

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
