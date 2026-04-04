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

package builtins

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"unicode"

	"github.com/cockroachdb/cockroach/pkg/sql/randgen"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/stretchr/testify/require"
)

func TestGeoBuiltinsInfo(t *testing.T) {
	defer leaktest.AfterTest(t)()

	for k, builtin := range geoBuiltins {
		t.Run(k, func(t *testing.T) {
			for i, overload := range builtin.overloads {
				t.Run(strconv.Itoa(i+1), func(t *testing.T) {
					infoFirstLine := strings.Trim(strings.Split(overload.Info, "\n\n")[0], "\t\n ")
					require.True(t, infoFirstLine[len(infoFirstLine)-1] == '.', "first line of info must end with a `.` character")
					require.True(t, unicode.IsUpper(rune(infoFirstLine[0])), "first character of info start with an uppercase letter.")
				})
			}
		})
	}
}

// TestGeoBuiltinsPointEmptyArgs tests POINT EMPTY arguments do not cause panics.
func TestGeoBuiltinsPointEmptyArgs(t *testing.T) {
	defer leaktest.AfterTest(t)()

	emptyGeometry, err := tree.ParseDGeometry("POINT EMPTY")
	require.NoError(t, err)
	emptyGeography, err := tree.ParseDGeography("POINT EMPTY")
	require.NoError(t, err)

	rng := rand.New(rand.NewSource(0))
	for k, builtin := range geoBuiltins {
		t.Run(k, func(t *testing.T) {
			for i, overload := range builtin.overloads {
				t.Run("overload_"+strconv.Itoa(i+1), func(t *testing.T) {
					for overloadIdx := 0; overloadIdx < overload.Types.Length(); overloadIdx++ {
						switch overload.Types.GetAt(overloadIdx).Family() {
						case types.GeometryFamily, types.GeographyFamily:
							t.Run("idx_"+strconv.Itoa(overloadIdx), func(t *testing.T) {
								var datums tree.Datums
								for i := 0; i < overload.Types.Length(); i++ {
									if i == overloadIdx {
										switch overload.Types.GetAt(i).Family() {
										case types.GeometryFamily:
											datums = append(datums, emptyGeometry)
										case types.GeographyFamily:
											datums = append(datums, emptyGeography)
										default:
											panic("unexpected condition")
										}
									} else {
										datums = append(datums, randgen.RandDatum(rng, overload.Types.GetAt(i), false))
									}
								}
								var call strings.Builder
								call.WriteString(k)
								call.WriteByte('(')
								for i, arg := range datums {
									if i > 0 {
										call.WriteString(", ")
									}
									call.WriteString(arg.String())
								}
								call.WriteByte(')')
								t.Logf("calling: %s", call.String())
								if overload.Fn != nil {
									_, _ = overload.Fn(&tree.EvalContext{}, datums)
								} else if overload.Generator != nil {
									_, _ = overload.Generator(&tree.EvalContext{}, datums)
								} else if overload.GeneratorWithExprs != nil {
									exprs := make(tree.Exprs, len(datums))
									for i := range datums {
										exprs[i] = datums[i]
									}
									_, _ = overload.GeneratorWithExprs(&tree.EvalContext{}, exprs)
								}
							})
						}
					}
				})
			}
		})
	}
}
