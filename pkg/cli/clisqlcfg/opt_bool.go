// Copyright 2021 The Cockroach Authors.
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

package clisqlcfg

import "strconv"

// OptBool implements a boolean value that may be undefined.
type OptBool struct {
	isDef bool
	v     bool
}

// Get returns whether the value is defined, and the value.
func (o OptBool) Get() (isDef, val bool) {
	return o.isDef, o.v
}

// String implements the pflag.Value interface.
func (o OptBool) String() string {
	if !o.isDef {
		return "<unspecified>"
	}
	return strconv.FormatBool(o.v)
}

// Type implements the pflag.Value interface.
func (o OptBool) Type() string { return "bool" }

// Set implements the pflag.Value interface.
func (o *OptBool) Set(v string) error {
	b, err := strconv.ParseBool(v)
	if err != nil {
		return err
	}
	o.isDef = true
	o.v = b
	return nil
}
