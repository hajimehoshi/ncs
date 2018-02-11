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
		In  string
		Out Color
	}{
		{
			In: "3000-Y10R",
			Out: Color{
				Blackness:     30,
				Chromaticness: 0,
				Hue:           10,
			},
		},
	}

	for _, c := range cases {
		got, err := Parse(c.In)
		if err != nil {
			t.Fatal(err)
		}
		want := c.Out
		if got != want {
			t.Errorf("Parse(%v): got: %v, want: %v", c.In, got, want)
		}
	}
}
