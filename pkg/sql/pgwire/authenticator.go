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

package pgwire

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/security"
)

// Authenticator is a component of an AuthMethod that determines if the
// given system identity (e.g.: Kerberos or X.509 principal, plain-old
// username, etc) is who it claims to be.
type Authenticator = func(
	ctx context.Context,
	systemIdentity security.SQLUsername,
	clientConnection bool,
	pwRetrieveFn PasswordRetrievalFn,
) error

// PasswordRetrievalFn defines a method to retrieve a hashed password
// and expiration time for a user logging in with password-based
// authentication.
type PasswordRetrievalFn = func(context.Context) (
	expired bool,
	pwHash security.PasswordHash,
	_ error,
)
