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

package colexecargs

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/testutils/buildutil"
)

func TestNoLinkForbidden(t *testing.T) {
	// Prohibit introducing any new dependencies into this package since it
	// should be very lightweight.
	buildutil.VerifyNoImports(t,
		"github.com/cockroachdb/cockroach/pkg/sql/colexec/colexecargs", true,
		nil /* forbiddenPkgs */, nil, /* forbiddenPrefixes */
		// allowlist:
		"github.com/cockroachdb/cockroach/pkg/col/coldata",
		"github.com/cockroachdb/cockroach/pkg/sql/colcontainer",
		"github.com/cockroachdb/cockroach/pkg/sql/colexecop",
		"github.com/cockroachdb/cockroach/pkg/sql/execinfra",
		"github.com/cockroachdb/cockroach/pkg/sql/execinfrapb",
		"github.com/cockroachdb/cockroach/pkg/sql/types",
		"github.com/cockroachdb/cockroach/pkg/util/mon",
		"github.com/marusama/semaphore",
	)
}
