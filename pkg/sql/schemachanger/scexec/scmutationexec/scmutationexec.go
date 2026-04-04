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

package scmutationexec

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scop"
)

// NewMutationVisitor creates a new scop.MutationVisitor.
func NewMutationVisitor(
	s MutationVisitorStateUpdater, nr NameResolver, sd SyntheticDescriptors, clock Clock,
) scop.MutationVisitor {
	return &visitor{
		nr:    nr,
		sd:    sd,
		s:     s,
		clock: clock,
	}
}

var _ scop.MutationVisitor = (*visitor)(nil)

type visitor struct {
	clock Clock
	nr    NameResolver
	sd    SyntheticDescriptors
	s     MutationVisitorStateUpdater
}

func (m *visitor) NotImplemented(_ context.Context, _ scop.NotImplemented) error {
	return nil
}
