// Copyright 2016 The Cockroach Authors.
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

package sql

import "github.com/cockroachdb/cockroach/pkg/sql/catalog/colinfo"

// planDataSource contains the data source information for data
// produced by a planNode.
type planDataSource struct {
	// columns gives the result columns (always anonymous source).
	columns colinfo.ResultColumns

	// plan which can be used to retrieve the data.
	plan planNode
}
