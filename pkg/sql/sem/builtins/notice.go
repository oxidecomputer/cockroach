// Copyright 2015 The Cockroach Authors.
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

package builtins

import (
	"strings"

	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgnotice"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/errors"
)

// crdbInternalSendNotice sends a notice.
// Note this is extracted to a different file to prevent churn on the pgwire
// test, which records line numbers.
func crdbInternalSendNotice(
	ctx *tree.EvalContext, severity string, msg string,
) (tree.Datum, error) {
	if ctx.ClientNoticeSender == nil {
		return nil, errors.AssertionFailedf("notice sender not set")
	}
	ctx.ClientNoticeSender.BufferClientNotice(
		ctx.Context,
		pgnotice.NewWithSeverityf(strings.ToUpper(severity), "%s", msg),
	)
	return tree.NewDInt(0), nil
}
