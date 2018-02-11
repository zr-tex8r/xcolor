// Copyright (c) 2017 Takayuki YATO (aka. "ZR")
//   GitHub:   https://github.com/zr-tex8r
//   Twitter:  @zr_tex8r
// Distributed under the MIT License.

package xcolor

import (
	"fmt"
	"image/color"
)

type rgbColor struct {
	r, g, b float64
}

func (*rgbColor) Model() Model {
	return Rgb
}

func (c *rgbColor) Params() []float64 {
	return []float64{c.r, c.g, c.b}
}

func (c *rgbColor) Compl() Color {
	return &rgbColor{1 - c.r, 1 - c.g, 1 - c.b}
}

func (c1 *rgbColor) Blend(r float64, c2 Color) Color {
	if r <= 0 {
		return c2
	} else if r >= 1 {
		return c1
	} else {
		c := c2.toRgb()
		return &rgbColor{bval(c1.r, c.r, r), bval(c1.g, c.g, r),
			bval(c1.b, c.b, r)}
	}
}

func (c *rgbColor) Convert(model Model) Color {
	switch model {
	case Gray:
		g := c.r*0.3 + c.g*0.59 + c.b*0.11
		return &grayColor{g}
	case Cmyk:
		return c.toCmyk()
	default:
		return c
	}
}

func (c *rgbColor) HtmlCode() string {
	return fmt.Sprintf("#%02X%02X%02X", vr8(c.r), vr8(c.g), vr8(c.b))
}

func (c *rgbColor) CssCode() string {
	return colorCssStr(c.Model(), c.Params(), stdPrec)
}

func (c *rgbColor) CssCodeWithPrec(prec int) string {
	return colorCssStr(c.Model(), c.Params(), prec)
}

func (c1 *rgbColor) toRgb() *rgbColor {
	return c1
}

func (c *rgbColor) toCmyk() *cmykColor {
	cc := &cmykColor{vr(1 - c.r), vr(1 - c.g), vr(1 - c.b), 0}
	return cc.normal()
}

func (c *rgbColor) String() string {
	return colorStr(c.Model(), c.Params())
}

func (c *rgbColor) GoColor() color.Color {
	return color.NRGBA{vr8(c.r), vr8(c.g), vr8(c.b), 0xFF}
}

func (c *rgbColor) RGBA() (r, g, b, a uint32) {
	return c.GoColor().RGBA()
}
