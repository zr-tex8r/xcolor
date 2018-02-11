// Copyright (c) 2017 Takayuki YATO (aka. "ZR")
//   GitHub:   https://github.com/zr-tex8r
//   Twitter:  @zr_tex8r
// Distributed under the MIT License.

package xcolor

import (
	"image/color"
	"math"
)

type cmykColor struct {
	c, m, y, k float64
}

func (*cmykColor) Model() Model {
	return Cmyk
}

func (c *cmykColor) Params() []float64 {
	return []float64{c.c, c.m, c.y, c.k}
}

func (c *cmykColor) normal() *cmykColor {
	dk := math.Min(c.c, math.Min(c.m, c.y))
	return &cmykColor{c.c - dk, c.m - dk, c.y - dk, c.k + dk}
}

func (c *cmykColor) denormal() *cmykColor {
	return &cmykColor{c.c + c.k, c.m + c.k, c.y + c.k, 0}
}

func (c *cmykColor) Compl() Color {
	c = c.denormal()
	c = &cmykColor{1 - c.c, 1 - c.m, 1 - c.y, 0}
	return c.normal()
}

func (c1 *cmykColor) Blend(r float64, c2 Color) Color {
	if r <= 0 {
		return c2
	} else if r >= 1 {
		return c1
	} else {
		c := c2.toCmyk()
		return &cmykColor{bval(c1.c, c.c, r), bval(c1.m, c.m, r),
			bval(c1.y, c.y, r), bval(c1.k, c.k, r)}
	}
}

func (c *cmykColor) Convert(model Model) Color {
	switch model {
	case Gray:
		return c.toRgb().Convert(model)
	case Rgb:
		return c.toRgb()
	default:
		return c
	}
}

func (c *cmykColor) HtmlCode() string {
	return c.toRgb().HtmlCode()
}

func (c *cmykColor) CssCode() string {
	return colorCssStr(c.Model(), c.Params(), stdPrec)
}

func (c *cmykColor) CssCodeWithPrec(prec int) string {
	return colorCssStr(c.Model(), c.Params(), prec)
}

func (c1 *cmykColor) toRgb() *rgbColor {
	c := c1.denormal()
	return &rgbColor{vr(1 - c.c), vr(1 - c.m), vr(1 - c.y)}
}

func (c1 *cmykColor) toCmyk() *cmykColor {
	return c1
}

func (c *cmykColor) String() string {
	return colorStr(c.Model(), c.Params())
}

func (c *cmykColor) GoColor() color.Color {
	return color.CMYK{vr8(c.c), vr8(c.m), vr8(c.y), vr8(c.k)}
}

func (c *cmykColor) RGBA() (r, g, b, a uint32) {
	return c.GoColor().RGBA()
}
