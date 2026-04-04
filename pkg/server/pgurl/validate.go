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

package pgurl

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/redact"
)

// Validate checks that the URL parameters are correct.
func (u *URL) Validate() error {
	var details bytes.Buffer
	var incorrect []redact.RedactableString

	switch u.net {
	case ProtoUnix:
		if !strings.HasPrefix(u.host, "/") {
			incorrect = append(incorrect, "host")
			fmt.Fprintln(&details, "Host parameter must start with '/' when using unix sockets.")
		}
		if u.sec != tnUnspecified && u.sec != tnNone {
			incorrect = append(incorrect, "sslmode")
			fmt.Fprintln(&details, "Cannot specify TLS settings when using unix sockets.")
		}
	case ProtoTCP:
		if strings.Contains(u.host, "/") {
			incorrect = append(incorrect, "host")
			fmt.Fprintln(&details, "Host parameter cannot contain '/' when using TCP.")
		}
	default:
		incorrect = append(incorrect, "net")
		fmt.Fprintln(&details, "Network protocol unspecified.")
	}

	if u.username == "" && u.hasPassword {
		incorrect = append(incorrect, "user")
		fmt.Fprintln(&details, "Username cannot be empty when a password is provided.")
	}

	switch u.authn {
	case authnClientCert, authnPasswordWithClientCert:
		if u.sec == tnUnspecified || u.sec == tnNone {
			incorrect = append(incorrect, "sslmode")
			fmt.Fprintln(&details, "Cannot use TLS client certificate authentication without a TLS transport.")
		}
		if u.clientCertPath == "" {
			incorrect = append(incorrect, "sslcert")
			fmt.Fprintln(&details, "Client certificate missing.")
		}
		if u.clientKeyPath == "" {
			incorrect = append(incorrect, "sslkey")
			fmt.Fprintln(&details, "Client key missing.")
		}
	case authnUndefined:
		incorrect = append(incorrect, "authn")
		fmt.Fprintln(&details, "Authentication method unspecified.")
	}

	if len(incorrect) > 0 {
		return errors.WithDetail(errors.Newf("URL validation error: %s", redact.Join(", ", incorrect)), details.String())
	}
	return nil
}
