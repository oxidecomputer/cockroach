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

// Package clierrorplus contains facilities that would nominally belong
// to package `clierror` instead, but which we do not wish to place
// there to prevent `clierror` from depending to more complex packages.
// We want `clierror` to remain lightweight so that the `cockroach-sql`
// standalone shell remains small.
package clierrorplus
