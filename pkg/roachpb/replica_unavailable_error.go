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
	context "context"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/gogo/protobuf/proto"
)

// NewReplicaUnavailableError initializes a new *ReplicaUnavailableError. It is
// provided with the range descriptor known to the replica, and the relevant
// replica descriptor within.
func NewReplicaUnavailableError(
	cause error, desc *RangeDescriptor, replDesc ReplicaDescriptor,
) error {
	return &ReplicaUnavailableError{
		Desc:    *desc,
		Replica: replDesc,
		Cause:   errors.EncodeError(context.Background(), cause),
	}
}

var _ errors.SafeFormatter = (*ReplicaUnavailableError)(nil)
var _ fmt.Formatter = (*ReplicaUnavailableError)(nil)
var _ errors.Wrapper = (*ReplicaUnavailableError)(nil)

// SafeFormatError implements errors.SafeFormatter.
func (e *ReplicaUnavailableError) SafeFormatError(p errors.Printer) error {
	p.Printf("replica unavailable: %s unable to serve request to %s: %s", e.Replica, e.Desc, e.Unwrap())
	return nil
}

// Format implements fmt.Formatter.
func (e *ReplicaUnavailableError) Format(s fmt.State, verb rune) { errors.FormatError(e, s, verb) }

// Error implements error.
func (e *ReplicaUnavailableError) Error() string {
	return fmt.Sprint(e)
}

// Unwrap implements errors.Wrapper.
func (e *ReplicaUnavailableError) Unwrap() error {
	return errors.DecodeError(context.Background(), e.Cause)
}

func init() {
	encode := func(ctx context.Context, err error) (msgPrefix string, safeDetails []string, payload proto.Message) {
		errors.As(err, &payload) // payload = err.(proto.Message)
		return "", nil, payload
	}
	decode := func(ctx context.Context, cause error, msgPrefix string, safeDetails []string, payload proto.Message) error {
		return payload.(*ReplicaUnavailableError)
	}
	typeName := errors.GetTypeKey((*ReplicaUnavailableError)(nil))
	errors.RegisterWrapperEncoder(typeName, encode)
	errors.RegisterWrapperDecoder(typeName, decode)
}
