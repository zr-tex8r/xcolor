// Copyright (c) 2017 Takayuki YATO (aka. "ZR")
//   GitHub:   https://github.com/zr-tex8r
//   Twitter:  @zr_tex8r
// Distributed under the MIT License.

package xcolor

// Enum Model represents a color model.
type Model int

const (
	Gray Model = iota
	Rgb
	Cmyk
)

var modelNames []string = []string{
	"gray", "rgb", "cmyk",
}

func (m Model) String() string {
	return modelNames[m]
}
