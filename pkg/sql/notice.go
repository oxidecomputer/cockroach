// Copyright 2020 The Cockroach Authors.
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

package sql

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/settings"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgnotice"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

// NoticesEnabled is the cluster setting that allows users
// to enable notices.
var NoticesEnabled = settings.RegisterBoolSetting(
	settings.TenantWritable,
	"sql.notices.enabled",
	"enable notices in the server/client protocol being sent",
	true,
).WithPublic()

// noticeSender is a subset of RestrictedCommandResult which allows
// sending notices.
type noticeSender interface {
	BufferNotice(pgnotice.Notice)
}

// BufferClientNotice implements the tree.ClientNoticeSender interface.
func (p *planner) BufferClientNotice(ctx context.Context, notice pgnotice.Notice) {
	if log.V(2) {
		log.Infof(ctx, "buffered notice: %+v", notice)
	}
	noticeSeverity, ok := pgnotice.ParseDisplaySeverity(pgerror.GetSeverity(notice))
	if !ok {
		noticeSeverity = pgnotice.DisplaySeverityNotice
	}
	if p.noticeSender == nil ||
		noticeSeverity > pgnotice.DisplaySeverity(p.SessionData().NoticeDisplaySeverity) ||
		!NoticesEnabled.Get(&p.execCfg.Settings.SV) {
		// Notice cannot flow to the client - because of one of these conditions:
		// * there is no client
		// * the session's NoticeDisplaySeverity is higher than the severity of the notice.
		// * the notice protocol was disabled
		return
	}
	p.noticeSender.BufferNotice(notice)
}
