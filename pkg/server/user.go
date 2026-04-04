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

	"github.com/cockroachdb/cockroach/pkg/server/serverpb"
	"github.com/cockroachdb/cockroach/pkg/sql/roleoption"
)

// UserSQLRoles return a list of the logged in SQL user roles.
func (s *baseStatusServer) UserSQLRoles(
	ctx context.Context, req *serverpb.UserSQLRolesRequest,
) (_ *serverpb.UserSQLRolesResponse, retErr error) {
	ctx = propagateGatewayMetadata(ctx)
	ctx = s.AnnotateCtx(ctx)

	username, isAdmin, err := s.privilegeChecker.getUserAndRole(ctx)
	if err != nil {
		return nil, err
	}

	var resp serverpb.UserSQLRolesResponse
	if !isAdmin {
		for name := range roleoption.ByName {
			hasRole, err := s.privilegeChecker.hasRoleOption(ctx, username, roleoption.ByName[name])
			if err != nil {
				return nil, err
			}
			if hasRole {
				resp.Roles = append(resp.Roles, name)
			}
		}
	} else {
		resp.Roles = append(resp.Roles, "ADMIN")
	}
	return &resp, nil
}
