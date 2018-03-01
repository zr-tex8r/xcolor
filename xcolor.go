// Copyright (c) 2018 Takayuki YATO (aka. "ZR")
//   GitHub:   https://github.com/zr-tex8r
//   Twitter:  @zr_tex8r
// Distributed under the MIT License.

// Package xcolor enables one to specify a color using color expression
// strings as used in the LaTeX package "xcolor".
//
// The package employs its own internal structs for representing colors,
// and mainly uses the structs. It also supports conversion to values
// of the standard color.Color type.
package xcolor

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

const minPrec int = 2
const stdPrec int = 3

func init() {
	changePrec(stdPrec)
	// assume that "current" color is always black
	repository["."] = NewGray(0)
	LoadPreset("default")
}

var repository = make(map[string]Color)

// LoadPreset loads a color name preset; the preset name "XXX" corresponds
// to the option "XXXnames" of the LaTeX xcolor package.
// The set names available are: "dvips", "svg".
func LoadPreset(name string) error {
	pre := preset[name]
	if pre == nil {
		return fmt.Errorf("Unknown preset name: %s", name)
	}
	for nam, col := range pre {
		repository[nam] = col
	}
	return nil
}

// GoColor converts a color expression to a color value of the standard
// color.Color type. Equivalent to Parse(expr).GoColor().
func GoColor(expr string) (clr color.Color, err error) {
	if xclr, err := Parse(expr); err == nil {
		return xclr.GoColor(), nil
	}
	return
}

// Parse converts a color expression to a Color value.
func Parse(expr string) (Color, error) {
	fs := strings.Split(expr, ":")
	switch len(fs) {
	case 1:
		return parseStdExpr(expr)
	case 2:
		return parseExtExpr(fs[0], fs[1])
	default:
		return nil, fmt.Errorf("Too many ':' in expression")
	}
}

func findColor(name string) (ret Color, err error) {
	name = strings.Trim(name, " ")
	if ret = repository[name]; ret != nil {
		return
	}
	return nil, fmt.Errorf("Unknown color name '%s'", name)
}

func findModel(name string) (Model, error) {
	name = strings.Trim(name, " ")
	for k, v := range modelNames {
		if v == name {
			return Model(k), nil
		}
	}
	return Model(0), fmt.Errorf("Unknown model name '%s'", name)
}

func parseNum(str string) (ret float64, err error) {
	str = strings.Trim(str, " ")
	if ret, err = strconv.ParseFloat(str, 64); err == nil {
		return
	}
	return ret, fmt.Errorf("Not a number ('%s')", str)
}

var white = NewGray(1)

func parseStdExpr(expr string) (ret Color, err error) {
	fs := strings.Split(expr, "!")
	if len(fs) > 2 && fs[len(fs)-2] == "" {
		return nil, fmt.Errorf("Postfix('!!') is not supported")
	}
	if ret, err = findColor(fs[0]); err != nil {
		return
	}
	for k := 1; k < len(fs); k += 2 {
		r, c := 0.0, white
		if r, err = parseNum(fs[k]); err != nil {
			return nil, err
		}
		if k+1 < len(fs) {
			if c, err = findColor(fs[k+1]); err != nil {
				return nil, err
			}
		}
		ret = ret.Blend(r/100, c)
	}
	return
}

func newFromParams(model Model, ps []float64) Color {
	switch model {
	case Rgb:
		return NewRgb(ps[0], ps[1], ps[2])
	case Cmyk:
		return NewCmyk(ps[0], ps[1], ps[2], ps[3])
	default:
		return NewGray(ps[0])
	}
}

func parseExtExpr(prefix, expr string) (ret Color, err error) {
	fs := strings.SplitN(prefix, ",", 2)
	model, err := findModel(fs[0])
	if err != nil {
		return
	}
	var div, cdiv float64
	if len(fs) > 1 {
		if div, err = parseNum(fs[1]); err != nil {
			return
		}
		if div <= 0 {
			return nil, fmt.Errorf("Negative div value (%v)", div)
		}
	}
	var param []float64
	for _, term := range strings.Split(expr, ";") {
		fs = strings.SplitN(term, ",", 2)
		if len(fs) < 2 {
			return nil, fmt.Errorf("Missing ',' in '%v'", term)
		}
		clr, err := parseStdExpr(fs[0])
		if err != nil {
			return ret, err
		}
		r, err := parseNum(fs[1])
		if err != nil {
			return ret, err
		}
		cdiv += r
		if param == nil {
			param = clr.Convert(model).Params()
			for k, v := range param {
				param[k] = v * r
			}
		} else {
			for k, v := range clr.Convert(model).Params() {
				param[k] += v * r
			}
		}
	}
	if div == 0 {
		div = cdiv
	}
	for k, _ := range param {
		param[k] /= div
	}
	ret = newFromParams(model, param)
	return
}
