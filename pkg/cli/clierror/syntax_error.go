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

package clierror

import (
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/errors"
	"github.com/lib/pq"
)

// IsSQLSyntaxError returns true iff the provided error is a SQL
// syntax error. The function works for the queries executed via the
// clisqlclient/clisqlexec packages.
func IsSQLSyntaxError(err error) bool {
	if pqErr := (*pq.Error)(nil); errors.As(err, &pqErr) {
		return string(pqErr.Code) == pgcode.Syntax.String()
	}
	return false
}
