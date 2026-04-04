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

package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/cockroachdb/cockroach/pkg/sql/types"
)

const (
	// aggKindTmplVar specifies the template "variable" that describes the kind
	// of aggregator using an aggregate function. It is replaced with "Hash",
	// "Ordered", or "Window" before executing the template.
	aggKindTmplVar = "_AGGKIND"
	hashAggKind    = "Hash"
	orderedAggKind = "Ordered"
	windowAggKind  = "Window"
)

func registerAggGenerator(aggGen generator, filenameSuffix, dep string, genWindowVariant bool) {
	aggGeneratorAdapter := func(aggKind string) generator {
		return func(inputFileContents string, wr io.Writer) error {
			inputFileContents = strings.ReplaceAll(inputFileContents, aggKindTmplVar, aggKind)
			return aggGen(inputFileContents, wr)
		}
	}
	aggKinds := []string{hashAggKind, orderedAggKind}
	if genWindowVariant {
		aggKinds = append(aggKinds, windowAggKind)
	}
	for _, aggKind := range aggKinds {
		registerGenerator(
			aggGeneratorAdapter(aggKind),
			fmt.Sprintf("%s_%s", strings.ToLower(aggKind), filenameSuffix),
			dep,
		)
	}
}

// aggTmplInfoBase is a helper struct used in generating the code of many
// aggregates serving as a base for calling methods on (whenever
// lastArgWidthOverload isn't available).
type aggTmplInfoBase struct {
	// canonicalTypeFamily is the canonical type family of the current aggregate
	// object used by the aggregate function.
	canonicalTypeFamily types.Family
}

// CopyVal is a function that should only be used in templates.
func (b *aggTmplInfoBase) CopyVal(dest, src string) string {
	return copyVal(b.canonicalTypeFamily, dest, src)
}

// SetVariableSize is a function that should only be used in templates. See the
// comment on setVariableSize for more details.
func (b aggTmplInfoBase) SetVariableSize(target, value string) string {
	return setVariableSize(b.canonicalTypeFamily, target, value)
}

// Remove unused warning.
var (
	a aggTmplInfoBase
	_ = a.SetVariableSize
	_ = a.CopyVal
)
