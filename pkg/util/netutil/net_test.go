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

package netutil

import (
	"fmt"
	"io"
	"net"
	"syscall"
	"testing"

	"github.com/cockroachdb/cmux"
	"github.com/cockroachdb/cockroach/pkg/util/contextutil"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestIsClosedConnection(t *testing.T) {
	for _, tc := range []struct {
		err           error
		isClosedError bool
	}{
		{
			fmt.Errorf("an error"),
			false,
		},
		{
			net.ErrClosed,
			true,
		},
		{
			cmux.ErrListenerClosed,
			true,
		},
		{
			grpc.ErrServerStopped,
			true,
		},
		{
			io.EOF,
			true,
		},
		{
			// TODO(rafi): should this be treated the same as EOF?
			io.ErrUnexpectedEOF,
			false,
		},
		{
			&net.AddrError{Err: "addr", Addr: "err"},
			true,
		},
		{
			syscall.ECONNRESET,
			true,
		},
		{
			syscall.EADDRINUSE,
			true,
		},
		{
			syscall.ECONNABORTED,
			true,
		},
		{
			syscall.ECONNREFUSED,
			true,
		},
		{
			syscall.EBADMSG,
			true,
		},
		{
			syscall.EINTR,
			false,
		},
		{
			&contextutil.TimeoutError{},
			false,
		},
	} {
		assert.Equalf(t, tc.isClosedError, IsClosedConnection(tc.err),
			"expected %q to be evaluated as %v", tc.err, tc.isClosedError,
		)
	}
}
