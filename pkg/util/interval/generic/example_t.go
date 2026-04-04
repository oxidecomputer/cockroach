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

package generic

import "github.com/cockroachdb/cockroach/pkg/roachpb"

//go:generate ./gen.sh *example generic

type example struct {
	id   uint64
	span roachpb.Span
}

// Methods required by generic contract.
func (ex *example) ID() uint64         { return ex.id }
func (ex *example) Key() []byte        { return ex.span.Key }
func (ex *example) EndKey() []byte     { return ex.span.EndKey }
func (ex *example) String() string     { return ex.span.String() }
func (ex *example) New() *example      { return new(example) }
func (ex *example) SetID(v uint64)     { ex.id = v }
func (ex *example) SetKey(v []byte)    { ex.span.Key = v }
func (ex *example) SetEndKey(v []byte) { ex.span.EndKey = v }
