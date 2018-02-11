// Copyright 2018 Hajime Hoshi
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

package ncs_test

import (
	"testing"

	. "github.com/hajimehoshi/ncs"
)

func TestParse(t *testing.T) {
	cases := []struct {
		In   string
		Out  Color
		Out2 string
	}{
		{
			In: "3010-Y10R",
			Out: Color{
				Blackness:     30,
				Chromaticness: 10,
				Hue:           10,
			},
			Out2: "3010-Y10R",
		},
		{
			In: "3030-Y10R",
			Out: Color{
				Blackness:     30,
				Chromaticness: 30,
				Hue:           10,
			},
			Out2: "3030-Y10R",
		},
		{
			In: "3080-Y10R",
			Out: Color{
				Blackness:     30,
				Chromaticness: 70,
				Hue:           10,
			},
			Out2: "3070-Y10R",
		},
		{
			In: "4020-B",
			Out: Color{
				Blackness:     40,
				Chromaticness: 20,
				Hue:           200,
			},
			Out2: "4020-B",
		},
		{
			In: "3020-B50G",
			Out: Color{
				Blackness:     30,
				Chromaticness: 20,
				Hue:           250,
			},
			Out2: "3020-B50G",
		},
		{
			In: "9910-B50G",
			Out: Color{
				Blackness:     99,
				Chromaticness: 1,
				Hue:           250,
			},
			Out2: "9901-B50G",
		},
		{
			In: "3000-N",
			Out: Color{
				Blackness:     30,
				Chromaticness: 0,
				Hue:           0,
			},
			Out2: "3000-N",
		},
		{
			// If Chromaticness is 0, Hue is ignored.
			In: "3000-Y10R",
			Out: Color{
				Blackness:     30,
				Chromaticness: 0,
				Hue:           0,
			},
			Out2: "3000-N",
		},
		{
			In: "3020-N",
			Out: Color{
				Blackness:     30,
				Chromaticness: 0,
				Hue:           0,
			},
			Out2: "3000-N",
		},
	}

	for _, c := range cases {
		got, err := Parse(c.In)
		if err != nil {
			t.Fatal(err)
		}
		want := c.Out
		if got != want {
			t.Errorf("Parse(%q): got %#v, want %#v", c.In, got, want)
		}

		got2 := c.Out2
		want2 := want.String()
		if got2 != want2 {
			t.Errorf("Parse(%q).String(): got %q, want %q", c.In, got2, want2)
		}
	}
}
