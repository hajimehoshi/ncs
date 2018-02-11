// Copyright 2017 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package ncs provides color in Natural Color System.
package ncs

import (
	"fmt"
	"regexp"
	"strconv"
)

type rgb struct {
	// Each value is in 0x00-0xff, but the type is uint16 for overflow.
	r uint16
	g uint16
	b uint16
}

var (
	// See https://en.wikipedia.org/wiki/Natural_Color_System
	yellow = &rgb{0xFF, 0xD3, 0x00}
	red    = &rgb{0xC4, 0x02, 0x33}
	blue   = &rgb{0x00, 0x87, 0xBD}
	green  = &rgb{0x00, 0x9F, 0x6B}
)

// Color represents a color in Natural Color System.
// Color implements image/color's Color interface.
//
// Color is well-formed when the below conditions are satisfied.
type Color struct {
	// 0 to 99
	Blackness int

	// 0 to min(100 - Blackness, 99)
	Chromaticness int

	// If Chromaticness is 0, Hue is 0.
	// Otherwise:
	//   0   to 99  for Y to R
	//   100 to 199 for R to B
	//   200 to 299 for B to G
	//   300 to 399 for G to Y
	Hue int
}

var (
	re = regexp.MustCompile(`\A(\d{2})(\d{2})-(N|Y|R|B|G|Y\d{2}R|R\d{2}B|B\d{2}G|G\d{2}Y)\z`)
)

// Parse parses a given string and returns a Color object.
// The conversion is approximate.
func Parse(str string) (Color, error) {
	m := re.FindStringSubmatch(str)
	if m == nil {
		return Color{}, fmt.Errorf("ncs: invalid format: %s", str)
	}
	b, err := strconv.Atoi(m[1])
	if err != nil {
		return Color{}, err
	}
	c, err := strconv.Atoi(m[2])
	if err != nil {
		return Color{}, err
	}
	h := 0
	switch {
	case m[3] == "N":
		c = 0
	case m[3] == "Y":
		h = 0
	case m[3] == "R":
		h = 100
	case m[3] == "B":
		h = 200
	case m[3] == "G":
		h = 300
	case m[3][0] == 'Y':
		h, err = strconv.Atoi(m[3][1:3])
	case m[3][0] == 'R':
		h, err = strconv.Atoi(m[3][1:3])
		h += 100
	case m[3][0] == 'B':
		h, err = strconv.Atoi(m[3][1:3])
		h += 200
	case m[3][0] == 'G':
		h, err = strconv.Atoi(m[3][1:3])
		h += 300
	default:
		panic("not reached")
	}
	if err != nil {
		panic("not reached")
	}
	if c > 100-b {
		c = 100 - b
	}
	if c == 0 {
		h = 0
	}
	return Color{
		Blackness:     b,
		Chromaticness: c,
		Hue:           h,
	}, nil
}

// String returns a string representing the color.
func (c Color) String() string {
	hue := "?"
	switch {
	case c.Chromaticness == 0:
		hue = "N"
	case c.Hue == 0:
		hue = "Y"
	case c.Hue == 100:
		hue = "R"
	case c.Hue == 200:
		hue = "B"
	case c.Hue == 300:
		hue = "G"
	case 0 < c.Hue && c.Hue < 100:
		hue = fmt.Sprintf("Y%02dR", c.Hue)
	case 100 < c.Hue && c.Hue < 200:
		hue = fmt.Sprintf("R%02dB", c.Hue-100)
	case 200 < c.Hue && c.Hue < 300:
		hue = fmt.Sprintf("B%02dG", c.Hue-200)
	case 300 < c.Hue && c.Hue < 400:
		hue = fmt.Sprintf("G%02dY", c.Hue-300)
	}
	return fmt.Sprintf("%02d%02d-%s", c.Blackness, c.Chromaticness, hue)
}

// RGBA implements Color's RGBA.
func (c Color) RGBA() (r, g, b, a uint32) {
	if c.Chromaticness == 0 {
		a := uint32(100-c.Blackness) * 0xffff / 100
		return a, a, a, 0xffff
	}
	var c0 *rgb
	var c1 *rgb
	v := uint16(0)
	switch {
	case 0 <= c.Hue && c.Hue < 100:
		c0 = yellow
		c1 = red
		v = uint16(c.Hue)
	case 100 <= c.Hue && c.Hue < 200:
		c0 = red
		c1 = blue
		v = uint16(c.Hue - 100)
	case 200 <= c.Hue && c.Hue < 300:
		c0 = blue
		c1 = green
		v = uint16(c.Hue - 200)
	case 300 <= c.Hue && c.Hue < 400:
		c0 = green
		c1 = yellow
		v = uint16(c.Hue - 300)
	}
	c2 := &rgb{
		r: (c0.r*(100-v) + c1.r*v) / 100,
		g: (c0.g*(100-v) + c1.g*v) / 100,
		b: (c0.b*(100-v) + c1.b*v) / 100,
	}
	ch := uint16(c.Chromaticness)
	cw := &rgb{
		r: (0xff*(100-ch) + c2.r*ch) / 100,
		g: (0xff*(100-ch) + c2.g*ch) / 100,
		b: (0xff*(100-ch) + c2.b*ch) / 100,
	}
	cb := &rgb{
		r: c2.r * ch / 100,
		g: c2.g * ch / 100,
		b: c2.b * ch / 100,
	}
	bl := uint16(c.Blackness)
	blmax := 100 - ch
	if blmax == 0 {
		return uint32(cw.r) * 0x101, uint32(cw.g) * 0x101, uint32(cw.b) * 0x101, 0xffff
	}
	if bl > blmax {
		return 0, 0, 0, 0xffff
	}
	c4 := &rgb{
		r: (cw.r*(blmax-bl) + cb.r*bl) / blmax,
		g: (cw.g*(blmax-bl) + cb.g*bl) / blmax,
		b: (cw.b*(blmax-bl) + cb.b*bl) / blmax,
	}
	return uint32(c4.r) * 0x101, uint32(c4.g) * 0x101, uint32(c4.b) * 0x101, 0xffff
}
