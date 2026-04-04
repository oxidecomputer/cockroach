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

package execinfra

import (
	"github.com/cockroachdb/cockroach/pkg/sql/rowenc/valueside"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
	"github.com/cockroachdb/errors"
)

// DecodeDatum decodes the given bytes slice into a datum of the given type. It
// returns an error if the decoding is not valid, or if there are any remaining
// bytes.
func DecodeDatum(datumAlloc *tree.DatumAlloc, typ *types.T, data []byte) (tree.Datum, error) {
	datum, rem, err := valueside.Decode(datumAlloc, typ, data)
	if err != nil {
		return nil, errors.NewAssertionErrorWithWrappedErrf(err,
			"error decoding %d bytes", errors.Safe(len(data)))
	}
	if len(rem) != 0 {
		return nil, errors.AssertionFailedf(
			"%d trailing bytes in encoded value", errors.Safe(len(rem)))
	}
	return datum, nil
}
