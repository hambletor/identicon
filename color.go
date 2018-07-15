package identicon

import (
	"image/color"
	"math"
)

// HSL represents Hue/Saturation/Luminace
type HSL struct {
	H float64
	S float64
	L float64
}

func complementary(c color.Color) color.Color {
	//convert color to hsl representation for easier manipulation
	hsl := colorToHSL(c)
	//complimentary color is 180 degrees away from original color
	switch {
	case hsl.H < 180:
		hsl.H = hsl.H + 180
	default:
		hsl.H = hsl.H - 180
	}
	//convert back to color.Color
	cmp := hslToRGB(*hsl)
	return cmp
}

func hslToRGB(hsl HSL) color.Color {
	c := (1 - math.Abs((2*hsl.L)-1)) * hsl.S
	y := (math.Mod(hsl.H/60, 2) - 1)
	x := c * (1 - math.Abs(float64(y)))
	m := hsl.L - (c / 2)
	var r, g, b float64
	switch {
	case 0 <= hsl.H && hsl.H <= 60:
		r = (c + m) * 255
		g = (x + m) * 255
		b = m * 255
	case 60 <= hsl.H && hsl.H < 120:
		r = (x + m) * 255
		g = (c + m) * 255
		b = m * 255
	case 120 <= hsl.H && hsl.H < 180:
		r = m * 255
		g = (c + m) * 255
		b = (x + m) * 255
	case 180 <= hsl.H && hsl.H < 240:
		r = m * 255
		g = (x + m) * 255
		b = (c + m) * 255
	case 240 <= hsl.H && hsl.H < 300:
		r = (x + m) * 255
		g = m * 255
		b = (c + m) * 255
	case 300 <= hsl.H && hsl.H < 360:
		r = (c + m) * 255
		g = m * 255
		b = (x + m) * 255
	}
	return color.RGBA{R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 255}
}

func colorToHSL(c color.Color) *HSL {
	hsl := HSL{}

	red, green, blue, _ := c.RGBA()
	r := float64(red/255) / float64(255.0)
	g := float64(green/255) / float64(255.0)
	b := float64(blue/255) / float64(255.0)
	min := math.Min(r, g)
	min = math.Min(min, b)
	max := math.Max(r, g)
	max = math.Max(max, b)
	delta := max - min

	// calculate Luminace
	hsl.L = ((min + max) / 2)

	// calculate Saturation
	switch {
	case delta != 0:
		hsl.S = delta / (1 - math.Abs(2*hsl.L-1))
	default:
		hsl.S = 0
	}

	// Calculate Hue
	switch {
	case delta == 0:
		hsl.H = 0
	case r == max: // If Red is max, then Hue = ((G-B)/(max-min)) * 60
		hsl.H = math.Mod((g-b)/delta, 6) * 60
	case g == max: // If Green is max, then Hue = 2.0 + ((B-R)/(max-min)) * 60
		hsl.H = (2 + ((b - r) / delta)) * 60
	case b == max: // If Blue is max, then Hue = 4.0 + ((R-G)/(max-min)) * 60
		hsl.H = (4 + ((r - g) / delta)) * 60
	}

	// if the Hue value is negative, add 360 degrees to make it positive
	if hsl.H < 0 {
		hsl.H = hsl.H + 360
	}
	return &hsl
}
