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

package errorutil

import "github.com/cockroachdb/cockroach/pkg/util/errorutil/unimplemented"

// UnsupportedWithMultiTenancy returns an error suitable for returning when an
// operation could not be carried out due to the SQL server running in
// multi-tenancy mode. In that mode, Gossip and other components of the KV layer
// are not available.
func UnsupportedWithMultiTenancy(issue int) error {
	return unimplemented.NewWithIssue(issue, "operation is unsupported in multi-tenancy mode")
}

// FeatureNotAvailableToNonSystemTenantsIssue is to be used with the
// Optional and related error interfaces when a feature is simply not
// available to non-system tenants (i.e. we're not planning to change
// this).
// For all other multitenancy errors where there is a plan to
// improve the situation, a specific issue should be created instead.
const FeatureNotAvailableToNonSystemTenantsIssue = 54252
