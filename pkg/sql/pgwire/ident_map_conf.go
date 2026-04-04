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
	"strings"

	"github.com/cockroachdb/cockroach/pkg/settings"
	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/identmap"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

// serverIdentityMapSetting is the name of the cluster setting that
// holds the pg_ident configuration.
const serverIdentityMapSetting = "server.identity_map.configuration"

var connIdentityMapConf = func() *settings.StringSetting {
	s := settings.RegisterValidatedStringSetting(
		settings.TenantWritable,
		serverIdentityMapSetting,
		"system-identity to database-username mappings",
		"",
		func(values *settings.Values, s string) error {
			_, err := identmap.From(strings.NewReader(s))
			return err
		},
	)
	s.SetVisibility(settings.Public)
	return s
}()

// loadLocalIdentityMapUponRemoteSettingChange initializes the local
// node's cache of the identity map configuration each time the cluster
// setting is updated.
func loadLocalIdentityMapUponRemoteSettingChange(
	ctx context.Context, server *Server, st *cluster.Settings,
) {
	val := connIdentityMapConf.Get(&st.SV)
	idMap, err := identmap.From(strings.NewReader(val))
	if err != nil {
		log.Ops.Warningf(ctx, "invalid %s: %v", serverIdentityMapSetting, err)
		idMap = identmap.Empty()
	}

	server.auth.Lock()
	defer server.auth.Unlock()
	server.auth.identityMap = idMap
}
