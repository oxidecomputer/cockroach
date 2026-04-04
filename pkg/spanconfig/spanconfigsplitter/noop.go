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

package spanconfigsplitter

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/spanconfig"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
)

var _ spanconfig.Splitter = &NoopSplitter{}

// NoopSplitter is a Splitter that only returns "illegal use" errors.
type NoopSplitter struct{}

// Splits is part of spanconfig.Splitter.
func (i NoopSplitter) Splits(context.Context, catalog.TableDescriptor) (int, error) {
	return 0, nil
}
