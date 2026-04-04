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

// Package tenantsettingswatcher implements an in-memory view of the
// tenant_settings table (containing overrides for tenant settings) using a
// rangefeed. This functionality is used on host cluster nodes, which allow
// tenants to retrieve the overrides and listen for changes.
package tenantsettingswatcher
