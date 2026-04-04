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

package server

import (
	"context"
	"net"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/errors"
)

// ListenError is returned from Start when we fail to start listening on either
// the main Cockroach port or the HTTP port, so that the CLI can instruct the
// user on what might have gone wrong.
type ListenError struct {
	cause error
	Addr  string
}

// Error implements error.
func (l *ListenError) Error() string { return l.cause.Error() }

// Unwrap is because ListenError is a wrapper.
func (l *ListenError) Unwrap() error { return l.cause }

// ListenAndUpdateAddrs starts a TCP listener on the specified address
// then updates the address and advertised address fields based on the
// actual interface address resolved by the OS during the Listen()
// call.
func ListenAndUpdateAddrs(
	ctx context.Context, addr, advertiseAddr *string, connName string,
) (net.Listener, error) {
	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		return nil, &ListenError{
			cause: err,
			Addr:  *addr,
		}
	}
	if err := base.UpdateAddrs(ctx, addr, advertiseAddr, ln.Addr()); err != nil {
		return nil, errors.Wrapf(err, "internal error: cannot parse %s listen address", connName)
	}
	return ln, nil
}
