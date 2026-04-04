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

package colexecbase

import "github.com/cockroachdb/cockroach/pkg/sql/sem/tree"

// BinaryOverloadHelper is a utility struct used for templates of the binary
// overloads that fall back to the row-based tree.Datum computation.
//
// In order for the templates to see it correctly, a local variable named
// `_overloadHelper` of this type must be declared before the inlined
// overloaded code.
type BinaryOverloadHelper struct {
	BinFn   tree.TwoArgFn
	EvalCtx *tree.EvalContext
}
