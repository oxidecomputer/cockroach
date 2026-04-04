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

package scop

import "github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"

//go:generate go run ./generate_visitor.go scop Validation validation.go validation_visitor_generated.go

type validationOp struct{ baseOp }

func (validationOp) Type() Type { return ValidationType }

// ValidateUniqueIndex validates uniqueness of entries for a unique index.
type ValidateUniqueIndex struct {
	validationOp
	TableID descpb.ID
	IndexID descpb.IndexID
}

// ValidateCheckConstraint validates a check constraint on a table's columns.
type ValidateCheckConstraint struct {
	validationOp
	TableID descpb.ID
	Name    string
}

// Make sure baseOp is used for linter.
var _ = validationOp{baseOp: baseOp{}}
