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

package validate

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/clusterversion"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
)

const InvalidSchemaChangerStatePrefix = invalidSchemaChangerStatePrefix + ":"

func TestingSchemaChangerState(
	ctx context.Context, desc catalog.Descriptor,
) catalog.ValidationErrors {
	vea := validationErrorAccumulator{
		targetLevel:       catalog.ValidationLevelSelfOnly,
		activeVersion:     clusterversion.TestingClusterVersion,
		currentState:      validatingDescriptor,
		currentDescriptor: desc,
	}
	validateSchemaChangerState(desc, &vea)
	return vea.errors
}
