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

package colexectestutils

import (
	"github.com/cockroachdb/cockroach/pkg/sql/colexec/colexecargs"
	"github.com/cockroachdb/cockroach/pkg/sql/colexecop"
)

// MakeInputs is a utility function that populates a slice of
// colexecargs.OpWithMetaInfo objects based on sources.
func MakeInputs(sources []colexecop.Operator) []colexecargs.OpWithMetaInfo {
	inputs := make([]colexecargs.OpWithMetaInfo, len(sources))
	for i := range sources {
		inputs[i].Root = sources[i]
	}
	return inputs
}
