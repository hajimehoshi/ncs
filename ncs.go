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
type Color struct {
	Blackness     int // 00 - 99
	Chromaticness int // 00 - (100 - Blackness)
	hue           int
	monochrome    bool // true is the hue is N
}

var (
	re = regexp.MustCompile(`\A(\d{2})(\d{2})-(N|Y|R|B|G|Y\d{2}R|R\d{2}B|B\d{2}G|G\d{2}Y)\z`)
)

// Parse parses a given string and returns a Color object.
// The conversion is approximate.
func Parse(str string) (*Color, error) {
	m := re.FindStringSubmatch(str)
	if m == nil {
		return nil, fmt.Errorf("ncs: invalid format: %s", str)
	}
	b, err := strconv.Atoi(m[1])
	if err != nil {
		return nil, err
	}
	c, err := strconv.Atoi(m[2])
	if err != nil {
		return nil, err
	}
	h := 0
	switch {
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
	return &Color{
		Blackness:     b,
		Chromaticness: c,
		monochrome:    m[3] == "N",
		hue:           h,
	}, nil
}

// RGBA implements Color's RGBA.
func (c *Color) RGBA() (r, g, b, a uint32) {
	if c.monochrome {
		a := uint32(100-c.Blackness) * 0xffff / 100
		return a, a, a, 0xffff
	}
	var c0 *rgb
	var c1 *rgb
	v := uint16(0)
	switch {
	case 0 <= c.hue && c.hue < 100:
		c0 = yellow
		c1 = red
		v = uint16(c.hue)
	case 100 <= c.hue && c.hue < 200:
		c0 = red
		c1 = blue
		v = uint16(c.hue - 100)
	case 200 <= c.hue && c.hue < 300:
		c0 = blue
		c1 = green
		v = uint16(c.hue - 200)
	case 300 <= c.hue && c.hue < 400:
		c0 = green
		c1 = yellow
		v = uint16(c.hue - 300)
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
