// Copyright 2014 The Cockroach Authors.
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

package addr

import (
	"net"
	"strings"

	"github.com/cockroachdb/errors"
)

// SplitHostPort is like net.SplitHostPort however it supports
// addresses without a port number. In that case, the provided port
// number is used.
func SplitHostPort(v string, defaultPort string) (addr string, port string, err error) {
	addr, port, err = net.SplitHostPort(v)
	if err != nil {
		var aerr *net.AddrError
		if errors.As(err, &aerr) {
			if strings.HasPrefix(aerr.Err, "too many colons") {
				// Maybe this was an IPv6 address using the deprecated syntax
				// without '[...]'? Try that to help the user with a hint.
				// Note: the following is valid even if defaultPort is empty.
				// (An empty port number is always a valid listen address.)
				maybeAddr := "[" + v + "]:" + defaultPort
				addr, port, err = net.SplitHostPort(maybeAddr)
				if err == nil {
					err = errors.WithHintf(
						errors.Newf("invalid address format: %q", v),
						"enclose IPv6 addresses within [...], e.g. \"[%s]\"", v)
				}
			} else if strings.HasPrefix(aerr.Err, "missing port") {
				// It's inconvenient that SplitHostPort doesn't know how to ignore
				// a missing port number. Oh well.
				addr, port, err = net.SplitHostPort(v + ":" + defaultPort)
			}
		}
	}
	if err == nil && port == "" {
		port = defaultPort
	}
	return addr, port, err
}
