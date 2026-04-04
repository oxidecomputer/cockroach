// Copyright 2022 The Cockroach Authors.
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

package poison

import (
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/errors"
)

// NewPoisonedError instantiates a *PoisonedError referencing a poisoned latch
// (as identified by span and timestamp).
func NewPoisonedError(span roachpb.Span, ts hlc.Timestamp) *PoisonedError {
	return &PoisonedError{Span: span, Timestamp: ts}
}

var _ errors.SafeFormatter = (*PoisonedError)(nil)
var _ fmt.Formatter = (*PoisonedError)(nil)

// SafeFormatError implements errors.SafeFormatter.
func (e *PoisonedError) SafeFormatError(p errors.Printer) error {
	p.Printf("encountered poisoned latch %s@%s", e.Span, e.Timestamp)
	return nil
}

// Format implements fmt.Formatter.
func (e *PoisonedError) Format(s fmt.State, verb rune) { errors.FormatError(e, s, verb) }

// Error implements error.
func (e *PoisonedError) Error() string {
	return fmt.Sprint(e)
}
