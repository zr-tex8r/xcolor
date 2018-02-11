// Copyright (c) 2017 Takayuki YATO (aka. "ZR")
//   GitHub:   https://github.com/zr-tex8r
//   Twitter:  @zr_tex8r
// Distributed under the MIT License.

package xcolor

import (
	"fmt"
	color "image/color"
	"strings"
)

// Interface Color represents a color.
//
// Note: although the type Color can be used as color.Color by itself, such
// use will be agnostic of the underlying color model, unless the program
// knows the xcolor package and treats its concrete types. On the other
// hand, the conversion using GoColor() will yield values of appropriate
// concrete types of 'image/color', namely,  NRGBA, CMYK and Gray.
type Color interface {
	// Model tells the model of this color.
	Model() Model
	// Params returns the values of the numerical parameters for this color.
	Params() []float64
	// Compl returns the complementary of this color.
	Compl() Color
	// Blend returns the blended color made of this and c.
	Blend(r float64, c Color) Color
	// Convert converts this color to another color model.
	Convert(m Model) Color
	// HtmlCode returns the HTML color literal for this color.
	HtmlCode() string
	// CssCode returns the CSS color literal for this color.
	CssCode() string
	// GoColor() convert to the Color value of standard color package.
	GoColor() color.Color
	// RGBA() is provided so that Color will satisfy color.Color.
	RGBA() (r, g, b, a uint32)

	toRgb() *rgbColor
	toCmyk() *cmykColor

	fmt.Stringer
}

func bval(p1, p2, r float64) float64 {
	r = vr(r)
	return p1*r + p2*(1-r)
}

func vr(p float64) float64 {
	if p <= 0 {
		return 0
	} else if p >= 1 {
		return 1
	} else {
		return p
	}
}

func vr8(p float64) uint8 {
	if p <= 0 {
		return 0
	} else if p >= 1 {
		return 0xFF
	} else {
		return uint8(p*0xFF + 0.5)
	}
}

// NewGray makes a gray color.
func NewGray(w float64) Color {
	return &grayColor{vr(w)}
}

// NewRgb makes an rgb color.
func NewRgb(r, g, b float64) Color {
	return &rgbColor{vr(r), vr(g), vr(b)}
}

// NewCmyk makes an rgb color.
func NewCmyk(c, m, y, k float64) Color {
	return &cmykColor{vr(c), vr(m), vr(y), vr(k)}
}

var prevPrec = -1
var paramFmt string
var pparamFmt string

func changePrec(prec int) {
	if prec < minPrec {
		prec = minPrec
	}
	if prec == prevPrec {
		return
	}
	prevPrec = prec
	paramFmt = fmt.Sprintf("%%.%df", prec)
	pparamFmt = fmt.Sprintf("%%.%df", prec-2)
}

func paramStr(v float64) string {
	s := fmt.Sprintf(paramFmt, v)
	return strings.TrimRight(strings.TrimRight(s, "0"), ".")
}

func pparamStr(v float64) string {
	s := fmt.Sprintf(pparamFmt, v*100)
	s = strings.TrimRight(strings.TrimRight(s, "0"), ".")
	return s + "%"
}

func colorStr(model Model, params []float64) string {
	changePrec(stdPrec)
	t := make([]string, len(params))
	for i := 0; i < len(params); i++ {
		t[i] = paramStr(params[i])
	}
	sps := strings.Join(t, ",")
	return fmt.Sprintf("[%v]{%v}", model, sps)
}

func colorCssStr(model Model, params []float64, prec int) string {
	changePrec(prec)
	t := make([]string, len(params))
	for i := 0; i < len(params); i++ {
		t[i] = pparamStr(vr(params[i]))
	}
	sps := strings.Join(t, ",")
	return fmt.Sprintf("%v(%v)", model, sps)
}
