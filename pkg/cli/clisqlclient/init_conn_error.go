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

package clisqlclient

import (
	"fmt"

	"github.com/cockroachdb/errors"
)

// InitialSQLConnectionError indicates that an error was encountered
// during the initial set-up of a SQL connection.
type InitialSQLConnectionError struct {
	err error
}

// Error implements the error interface.
func (i *InitialSQLConnectionError) Error() string { return i.err.Error() }

// Cause implements causer for compatibility with pkg/errors.
// NB: this is obsolete. Use Unwrap() instead.
func (i *InitialSQLConnectionError) Cause() error { return i.err }

// Unwrap implements errors.Wrapper.
func (i *InitialSQLConnectionError) Unwrap() error { return i.err }

// Format implements fmt.Formatter.
func (i *InitialSQLConnectionError) Format(s fmt.State, verb rune) { errors.FormatError(i, s, verb) }

// FormatError implements errors.Formatter.
func (i *InitialSQLConnectionError) FormatError(p errors.Printer) error {
	if p.Detail() {
		p.Print("error while establishing the SQL session")
	}
	return i.err
}
