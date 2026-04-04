// Copyright 2018 The Cockroach Authors.
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

package cat

import (
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

// DataSourceName is an alias for tree.TableName, and is used for views and
// sequences as well as tables.
type DataSourceName = tree.TableName

// DataSource is an interface to a database object that provides rows, like a
// table, a view, or a sequence.
type DataSource interface {
	Object

	// Name returns the unqualified name of the object.
	Name() tree.Name

	// CollectTypes returns all user defined types that the column uses.
	// This includes types used in default expressions, computed columns,
	// and the type of the column itself.
	CollectTypes(ord int) (descpb.IDs, error)
}
