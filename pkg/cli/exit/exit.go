// Copyright 2020 The Cockroach Authors.
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

package exit

import (
	"fmt"
	"os"

	"github.com/cockroachdb/redact"
)

// Code represents an exit code.
type Code struct {
	code int
}

// String implements the fmt.Stringer interface.
func (c Code) String() string { return fmt.Sprint(c.code) }

// Format implements the fmt.Formatter interface.
func (c Code) Format(s fmt.State, verb rune) {
	_, f := redact.MakeFormat(s, verb)
	fmt.Fprintf(s, f, c.code)
}

// SafeValue implements the redact.SafeValue interface.
func (c Code) SafeValue() {}

var _ redact.SafeValue = Code{}

// WithCode terminates the process and sets its exit status code to
// the provided code.
func WithCode(code Code) {
	os.Exit(code.code)
}
