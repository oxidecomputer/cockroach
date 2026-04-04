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

package paramparse

import (
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

// UniqueConstraintParamContext is used as a validation context for
// ValidateUniqueConstraintParams function.
// IsPrimaryKey: set to true if the unique constraint for primary key.
// IsSharded: set to true if the unique constraint has a hash sharded index.
type UniqueConstraintParamContext struct {
	IsPrimaryKey bool
	IsSharded    bool
}

// ValidateUniqueConstraintParams checks if there is any storage parameters
// invalid as a param for Unique Constraint.
func ValidateUniqueConstraintParams(
	params tree.StorageParams, ctx UniqueConstraintParamContext,
) error {
	// Only `bucket_count` is allowed for primary key and unique index.
	for _, param := range params {
		switch param.Key {
		case `bucket_count`:
			if ctx.IsSharded {
				continue
			}
			return pgerror.New(
				pgcode.InvalidParameterValue,
				`"bucket_count" storage param should only be set with "USING HASH" for hash sharded index`,
			)
		default:
			if ctx.IsPrimaryKey {
				return pgerror.Newf(pgcode.InvalidParameterValue, "invalid storage param %q on primary key", params[0].Key)
			}
			return pgerror.Newf(pgcode.InvalidParameterValue, "invalid storage param %q on unique index", params[0].Key)
		}
	}
	return nil
}
