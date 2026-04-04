// Copyright 2019 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package reduce

type intPass struct {
	name string
	fn   func(string, int) (string, bool, error)
}

// MakeIntPass returns a Pass with a state that starts at 0 and increments by
// 1 each Advance. f is a transformation function that takes the input and an
// index (zero-based) determining which occurrence to transform. It returns the
// possibly transformed output and a boolean that is false if a transformation
// could not be done because i was exhausted.
//
// For example, if f is a func that replaces spaces with underscores, invoking
// that function with an i value of 2 should return the input string with only
// the 3rd occurrence of a space replaced with an underscore.
//
// This is a convenience wrapper since a large number of Pass implementations
// just need their state to increment a counter and don't have to keep track of
// other things like byte offsets.
func MakeIntPass(name string, f func(s string, i int) (out string, ok bool, err error)) Pass {
	return intPass{
		name: name,
		fn:   f,
	}
}

func (p intPass) Name() string {
	return p.name
}

func (p intPass) New(string) State {
	return 0
}

func (p intPass) Transform(f string, s State) (string, Result, error) {
	i := s.(int)
	data, ok, err := p.fn(f, i)
	res := OK
	if !ok {
		res = STOP
	}
	return data, res, err
}

func (p intPass) Advance(f string, s State) State {
	return s.(int) + 1
}
