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

package migration

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/clusterversion"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

// FenceVersionFor constructs the appropriate "fence version" for the given
// cluster version. Fence versions allow the migrations infrastructure to safely
// step through consecutive cluster versions in the presence of Nodes (running
// any binary version) being added to the cluster. See the migration manager
// above for intended usage.
//
// Fence versions (and the migrations infrastructure entirely) were introduced
// in the 21.1 release cycle. In the same release cycle, we introduced the
// invariant that new user-defined versions (users being crdb engineers) must
// always have even-numbered Internal versions, thus reserving the odd numbers
// to slot in fence versions for each cluster version. See top-level
// documentation in pkg/clusterversion for more details.
func FenceVersionFor(
	ctx context.Context, cv clusterversion.ClusterVersion,
) clusterversion.ClusterVersion {
	if (cv.Internal % 2) != 0 {
		log.Fatalf(ctx, "only even numbered internal versions allowed, found %s", cv.Version)
	}

	// We'll pick the odd internal version preceding the cluster version,
	// slotting ourselves right before it.
	fenceCV := cv
	fenceCV.Internal--
	return fenceCV
}
