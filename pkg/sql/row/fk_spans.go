// Copyright 2019 The Cockroach Authors.
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

package row

import (
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/span"
)

// FKCheckSpan returns a span that can be scanned to ascertain existence of a
// specific row in a given index.
func FKCheckSpan(
	builder *span.Builder,
	splitter span.Splitter,
	values []tree.Datum,
	colMap catalog.TableColMap,
	numCols int,
) (roachpb.Span, error) {
	span, containsNull, err := builder.SpanFromDatumRow(values, numCols, colMap)
	if err != nil {
		return roachpb.Span{}, err
	}
	return splitter.ExistenceCheckSpan(span, numCols, containsNull), nil
}
