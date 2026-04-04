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

package reduce_test

import (
	"context"
	"go/parser"
	"regexp"
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/cmd/reduce/reduce"
)

func TestReduceGo(t *testing.T) {
	reduce.Walk(t, "testdata", nil /* filter */, isInterestingGo, reduce.ModeInteresting,
		nil /* chunkReducer */, goPasses)
}

var (
	goPasses = []reduce.Pass{
		removeLine,
		simplifyConsts,
	}
	removeLine = reduce.MakeIntPass("remove line", func(s string, i int) (string, bool, error) {
		sp := strings.Split(s, "\n")
		if i >= len(sp) {
			return "", false, nil
		}
		out := strings.Join(append(sp[:i], sp[i+1:]...), "\n")
		return out, true, nil
	})
	simplifyConstsRE = regexp.MustCompile(`[a-z0-9][a-z0-9]+`)
	simplifyConsts   = reduce.MakeIntPass("simplify consts", func(s string, i int) (string, bool, error) {
		out := simplifyConstsRE.ReplaceAllStringFunc(s, func(found string) string {
			i--
			if i == -1 {
				return found[:1]
			}
			return found
		})
		return out, i < 0, nil
	})
)

func isInterestingGo(contains string) reduce.InterestingFn {
	return func(ctx context.Context, f string) (bool, func()) {
		_, err := parser.ParseExpr(f)
		if err == nil {
			return false, nil
		}
		return strings.Contains(err.Error(), contains), nil
	}
}
