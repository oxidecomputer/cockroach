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

package roachpb

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/gogo/protobuf/proto"
)

// NewAmbiguousResultErrorf initializes a new AmbiguousResultError with
// an explanatory format and set of arguments.
func NewAmbiguousResultErrorf(format string, args ...interface{}) *AmbiguousResultError {
	return NewAmbiguousResultError(errors.NewWithDepthf(1, format, args...))
}

// NewAmbiguousResultError returns an AmbiguousResultError wrapping (via
// errors.Wrapper) the supplied error.
func NewAmbiguousResultError(err error) *AmbiguousResultError {
	return &AmbiguousResultError{
		EncodedError:      errors.EncodeError(context.Background(), err),
		DeprecatedMessage: err.Error(),
	}
}

var _ errors.SafeFormatter = (*AmbiguousResultError)(nil)
var _ fmt.Formatter = (*AmbiguousResultError)(nil)
var _ errors.Wrapper = func() errors.Wrapper {
	aErr := (*AmbiguousResultError)(nil)
	typeKey := errors.GetTypeKey(aErr)
	errors.RegisterWrapperEncoder(typeKey, func(ctx context.Context, err error) (msgPrefix string, safeDetails []string, payload proto.Message) {
		errors.As(err, &payload)
		return "", nil, payload
	})
	errors.RegisterWrapperDecoder(typeKey, func(ctx context.Context, cause error, msgPrefix string, safeDetails []string, payload proto.Message) error {
		return payload.(*AmbiguousResultError)
	})

	return aErr
}()

// SafeFormatError implements errors.SafeFormatter.
func (e *AmbiguousResultError) SafeFormatError(p errors.Printer) error {
	p.Printf("result is ambiguous: %s", e.unwrapOrDefault())
	return nil
}

// Format implements fmt.Formatter.
func (e *AmbiguousResultError) Format(s fmt.State, verb rune) { errors.FormatError(e, s, verb) }

// Error implements error.
func (e *AmbiguousResultError) Error() string {
	return fmt.Sprint(e)
}

// Unwrap implements errors.Wrapper.
func (e *AmbiguousResultError) Unwrap() error {
	if e.EncodedError.Error == nil {
		return nil
	}
	return errors.DecodeError(context.Background(), e.EncodedError)
}

func (e *AmbiguousResultError) unwrapOrDefault() error {
	cause := e.Unwrap()
	if cause == nil {
		return errors.New("unknown cause") // can be removed in 22.2
	}
	return cause
}

func (e *AmbiguousResultError) message(_ *Error) string {
	return fmt.Sprintf("result is ambiguous: %v", e.unwrapOrDefault())
}

// Type is part of the ErrorDetailInterface.
func (e *AmbiguousResultError) Type() ErrorDetailType {
	return AmbiguousResultErrType
}
