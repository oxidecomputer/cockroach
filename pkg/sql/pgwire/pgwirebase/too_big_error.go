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

package pgwirebase

import (
	"strconv"

	"github.com/cockroachdb/errors"
)

// withMessageTooBig decorates an error when a read would overflow the ReadBuffer.
type withMessageTooBig struct {
	cause error
	size  int
}

var _ error = (*withMessageTooBig)(nil)
var _ errors.SafeDetailer = (*withMessageTooBig)(nil)

func (w *withMessageTooBig) Error() string         { return w.cause.Error() }
func (w *withMessageTooBig) Unwrap() error         { return w.cause }
func (w *withMessageTooBig) SafeDetails() []string { return []string{strconv.Itoa(w.size)} }

// withMessageTooBigError decorates the error with a severity.
func withMessageTooBigError(err error, size int) error {
	if err == nil {
		return nil
	}

	return &withMessageTooBig{cause: err, size: size}
}

// IsMessageTooBigError denotes whether a message is too big.
func IsMessageTooBigError(err error) bool {
	var c withMessageTooBig
	return errors.HasType(err, &c)
}

// GetMessageTooBigSize attempts to unwrap and find a MessageTooBig.
func GetMessageTooBigSize(err error) int {
	if c := (*withMessageTooBig)(nil); errors.As(err, &c) {
		return c.size
	}
	return -1
}
