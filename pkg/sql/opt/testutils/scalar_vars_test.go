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

package testutils

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/opt"
	"github.com/stretchr/testify/assert"
)

func TestScalarVars(t *testing.T) {
	var md opt.Metadata
	var sv ScalarVars

	// toStr recreates the variable definitions from md and ScalarVars.
	toStr := func() string {
		var buf bytes.Buffer
		for i := 0; i < md.NumColumns(); i++ {
			id := opt.ColumnID(i + 1)
			m := md.ColumnMeta(id)
			if i > 0 {
				buf.WriteString(", ")
			}
			fmt.Fprintf(&buf, "%s %s", m.Alias, m.Type)
			if sv.NotNullCols().Contains(id) {
				buf.WriteString(" not null")
			}
		}
		return buf.String()
	}

	vars := "a int, b string not null, c decimal"
	md.Init()
	if err := sv.Init(&md, strings.Split(vars, ", ")); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, toStr(), vars)
}
