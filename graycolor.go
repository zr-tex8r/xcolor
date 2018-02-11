// Copyright (c) 2017 Takayuki YATO (aka. "ZR")
//   GitHub:   https://github.com/zr-tex8r
//   Twitter:  @zr_tex8r
// Distributed under the MIT License.

package xcolor

import (
	"image/color"
)

type grayColor struct {
	gray float64
}

func (*grayColor) Model() Model {
	return Gray
}

func (c *grayColor) Params() []float64 {
	return []float64{c.gray}
}

func (c *grayColor) Compl() Color {
	return &grayColor{1 - c.gray}
}

func (c1 *grayColor) Blend(r float64, c2 Color) Color {
	if r <= 0 {
		return c2
	} else if r >= 1 {
		return c1
	} else if c2.Model() == Gray {
		c := c2.(*grayColor)
		return &grayColor{bval(c1.gray, c.gray, r)}
	} else {
		// choose a non-gray color model
		return c2.Blend(1-r, c1)
	}
}

func (c *grayColor) Convert(model Model) Color {
	switch model {
	case Rgb:
		return c.toRgb()
	case Cmyk:
		return c.toCmyk()
	default:
		return c
	}
}

func (c *grayColor) HtmlCode() string {
	return c.toRgb().HtmlCode()
}

func (c *grayColor) CssCode() string {
	return colorCssStr(c.Model(), c.Params(), stdPrec)
}

func (c *grayColor) CssCodeWithPrec(prec int) string {
	return colorCssStr(c.Model(), c.Params(), prec)
}

func (c *grayColor) toRgb() *rgbColor {
	return &rgbColor{c.gray, c.gray, c.gray}
}

func (c *grayColor) toCmyk() *cmykColor {
	return &cmykColor{0, 0, 0, 1 - c.gray}
}

func (c *grayColor) String() string {
	return colorStr(c.Model(), c.Params())
}

func (c *grayColor) GoColor() color.Color {
	return color.Gray{vr8(c.gray)}
}

func (c *grayColor) RGBA() (r, g, b, a uint32) {
	return c.GoColor().RGBA()
}
