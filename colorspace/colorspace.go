// Package colorspace converts colors between sRGB, CIE XYZ, CIE L*a*b*,
// and HSL, using the D65 white point throughout so a value can be pushed
// through any chain of conversions without a whitepoint mismatch creeping
// in.
package colorspace

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// RGB is a color in the sRGB space, gamma-encoded, components in [0,1].
type RGB struct {
	R, G, B float64
}

// XYZ is a color in CIE 1931 XYZ space, relative to the D65 white point.
type XYZ struct {
	X, Y, Z float64
}

// Lab is a color in CIE L*a*b* space, relative to the D65 white point.
type Lab struct {
	L, A, B float64
}

// HSL is a color in hue/saturation/lightness space. H is in degrees
// [0,360); S and L are in [0,1].
type HSL struct {
	H, S, L float64
}

// D65 reference white, 2-degree observer.
const (
	whiteX = 0.95047
	whiteY = 1.00000
	whiteZ = 1.08883
)

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

// srgbToLinear and linearToSRGB implement the piecewise sRGB transfer
// function (IEC 61966-2-1), not a plain gamma curve.
func srgbToLinear(c float64) float64 {
	if c <= 0.04045 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, 2.4)
}

func linearToSRGB(c float64) float64 {
	if c <= 0.0031308 {
		return c * 12.92
	}
	return 1.055*math.Pow(c, 1/2.4) - 0.055
}

// ToXYZ converts an sRGB color to CIE XYZ (D65) via linear RGB and the
// standard sRGB primaries matrix.
func (c RGB) ToXYZ() XYZ {
	r := srgbToLinear(c.R)
	g := srgbToLinear(c.G)
	b := srgbToLinear(c.B)
	return XYZ{
		X: r*0.4124564 + g*0.3575761 + b*0.1804375,
		Y: r*0.2126729 + g*0.7151522 + b*0.0721750,
		Z: r*0.0193339 + g*0.1191920 + b*0.9503041,
	}
}

// ToRGB converts a CIE XYZ (D65) color to sRGB, clamping components that
// fall outside the sRGB gamut.
func (c XYZ) ToRGB() RGB {
	r := c.X*3.2404542 + c.Y*-1.5371385 + c.Z*-0.4985314
	g := c.X*-0.9692660 + c.Y*1.8760108 + c.Z*0.0415560
	b := c.X*0.0556434 + c.Y*-0.2040259 + c.Z*1.0572252
	return RGB{
		R: clamp01(linearToSRGB(r)),
		G: clamp01(linearToSRGB(g)),
		B: clamp01(linearToSRGB(b)),
	}
}

func labF(t float64) float64 {
	const delta = 6.0 / 29.0
	if t > delta*delta*delta {
		return math.Cbrt(t)
	}
	return t/(3*delta*delta) + 4.0/29.0
}

func labFInv(t float64) float64 {
	const delta = 6.0 / 29.0
	if t > delta {
		return t * t * t
	}
	return 3 * delta * delta * (t - 4.0/29.0)
}

// ToLab converts a CIE XYZ (D65) color to CIE L*a*b*.
func (c XYZ) ToLab() Lab {
	fx := labF(c.X / whiteX)
	fy := labF(c.Y / whiteY)
	fz := labF(c.Z / whiteZ)
	return Lab{
		L: 116*fy - 16,
		A: 500 * (fx - fy),
		B: 200 * (fy - fz),
	}
}

// ToXYZ converts a CIE L*a*b* color back to CIE XYZ (D65).
func (c Lab) ToXYZ() XYZ {
	fy := (c.L + 16) / 116
	fx := fy + c.A/500
	fz := fy - c.B/200
	return XYZ{
		X: whiteX * labFInv(fx),
		Y: whiteY * labFInv(fy),
		Z: whiteZ * labFInv(fz),
	}
}

// ToLab is a convenience wrapper for RGB -> XYZ -> Lab.
func (c RGB) ToLab() Lab { return c.ToXYZ().ToLab() }

// ToRGB is a convenience wrapper for Lab -> XYZ -> RGB.
func (c Lab) ToRGB() RGB { return c.ToXYZ().ToRGB() }

// ToHSL converts an sRGB color to HSL.
func (c RGB) ToHSL() HSL {
	max := math.Max(c.R, math.Max(c.G, c.B))
	min := math.Min(c.R, math.Min(c.G, c.B))
	l := (max + min) / 2

	if max == min {
		return HSL{H: 0, S: 0, L: l}
	}

	d := max - min
	var s float64
	if l > 0.5 {
		s = d / (2 - max - min)
	} else {
		s = d / (max + min)
	}

	var h float64
	switch max {
	case c.R:
		h = (c.G - c.B) / d
		if c.G < c.B {
			h += 6
		}
	case c.G:
		h = (c.B-c.R)/d + 2
	default:
		h = (c.R-c.G)/d + 4
	}
	h *= 60

	return HSL{H: h, S: s, L: l}
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t++
	}
	if t > 1 {
		t--
	}
	switch {
	case t < 1.0/6:
		return p + (q-p)*6*t
	case t < 1.0/2:
		return q
	case t < 2.0/3:
		return p + (q-p)*(2.0/3-t)*6
	default:
		return p
	}
}

// ToRGB converts an HSL color to sRGB.
func (c HSL) ToRGB() RGB {
	if c.S == 0 {
		return RGB{R: c.L, G: c.L, B: c.L}
	}

	var q float64
	if c.L < 0.5 {
		q = c.L * (1 + c.S)
	} else {
		q = c.L + c.S - c.L*c.S
	}
	p := 2*c.L - q
	h := c.H / 360

	return RGB{
		R: hueToRGB(p, q, h+1.0/3),
		G: hueToRGB(p, q, h),
		B: hueToRGB(p, q, h-1.0/3),
	}
}

// ParseHex parses a "#rrggbb" or "rrggbb" string into an RGB value.
func ParseHex(s string) (RGB, error) {
	s = strings.TrimPrefix(s, "#")
	if len(s) != 6 {
		return RGB{}, fmt.Errorf("colorspace: hex color must be 6 digits, got %q", s)
	}
	v, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return RGB{}, fmt.Errorf("colorspace: invalid hex color %q: %w", s, err)
	}
	return RGB{
		R: float64((v>>16)&0xff) / 255,
		G: float64((v>>8)&0xff) / 255,
		B: float64(v&0xff) / 255,
	}, nil
}

// Hex renders an RGB value as a "#rrggbb" string, rounding and clamping
// each channel to a byte.
func (c RGB) Hex() string {
	r := int(math.Round(clamp01(c.R) * 255))
	g := int(math.Round(clamp01(c.G) * 255))
	b := int(math.Round(clamp01(c.B) * 255))
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}
