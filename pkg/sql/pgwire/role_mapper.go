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
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/hba"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/identmap"
	"github.com/cockroachdb/errors"
)

// RoleMapper defines a mechanism by which an AuthMethod associated
// with an incoming connection may replace the caller-provided system
// identity (e.g.: GSSAPI or X.509 principal, LDAP DN, etc.) with zero
// or more SQLUsernames that will be subsequently validated against the
// SQL roles defined within the database. The mapping from system
// identity to database roles may be derived from the host-based
// authentication mechanism built into CockroachDB, or it could
// conceivably be implemented by an external directory service which
// maps groups of users onto database roles.
type RoleMapper = func(
	ctx context.Context,
	systemIdentity security.SQLUsername,
) ([]security.SQLUsername, error)

// UseProvidedIdentity is a trivial implementation of RoleMapper which always
// returns its input.
func UseProvidedIdentity(
	_ context.Context, id security.SQLUsername,
) ([]security.SQLUsername, error) {
	return []security.SQLUsername{id}, nil
}

var _ RoleMapper = UseProvidedIdentity

// HbaMapper implements the "map" option that may be defined in a
// host-based authentication rule. If the HBA entry does not define a
// "map" option, this function will return UseProvidedIdentity.
//
// This mapper will return an error if an applied mapping rule results
// in the root user or a reserved user, which includes the node,
// "public", and various other magic prefixes.
func HbaMapper(hbaEntry *hba.Entry, identMap *identmap.Conf) RoleMapper {
	mapName := hbaEntry.GetOption("map")
	if mapName == "" {
		return UseProvidedIdentity
	}
	return func(_ context.Context, id security.SQLUsername) ([]security.SQLUsername, error) {
		users, err := identMap.Map(mapName, id.Normalized())
		if err != nil {
			return nil, err
		}
		for _, user := range users {
			if user.IsRootUser() || user.IsReserved() {
				return nil, errors.Newf("system identity %q mapped to reserved database role %q",
					id.Normalized(), user.Normalized())
			}
		}
		return users, nil
	}
}
