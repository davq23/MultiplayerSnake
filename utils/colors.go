package utils

import (
	"math"
	"strconv"
	"strings"
)

// GenerateRandomColors generates random colors given an HSV space
func GenerateRandomColors(h float64, s float64, v float64) (int64, int64, int64) {
	hI := h * 6
	f := h*6 - hI
	p := v * (1 - s)
	q := v * (1 - f*s)
	t := v * (1 - (1-f)*s)

	var r, g, b float64

	switch int64(hI) {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	case 5:
		r, g, b = v, p, q
	}

	return int64(r * 256), int64(g * 256), int64(b * 256)
}

// RGB2Hex Converts RGB to Hex string
func RGB2Hex(r int64, g int64, b int64) string {
	hex := make([]string, 7)

	hex[0] = "#"

	rDivided, gDivided, bDivided := float64(r)/16, float64(g)/16, float64(b)/16

	hex[1] = strconv.FormatInt(int64(rDivided), 16)
	hex[2] = strconv.FormatInt(int64(math.Abs(rDivided-float64(int64(rDivided)))*16), 16)
	hex[3] = strconv.FormatInt(int64(gDivided), 16)
	hex[4] = strconv.FormatInt(int64(math.Abs(gDivided-float64(int64(gDivided)))*16), 16)
	hex[5] = strconv.FormatInt(int64(float64(b)/16), 16)
	hex[6] = strconv.FormatInt(int64(math.Abs(bDivided-float64(int64(bDivided)))*16), 16)

	return strings.Join(hex, "")
}
