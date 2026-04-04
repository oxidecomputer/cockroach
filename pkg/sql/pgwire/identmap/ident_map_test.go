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

package identmap

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentityMapElement(t *testing.T) {
	exactMatch := func(sysName, dbname string) element {
		return element{
			dbUser:       dbname,
			pattern:      regexp.MustCompile("^" + regexp.QuoteMeta(sysName) + "$"),
			substituteAt: -1,
		}
	}
	regexMatch := func(sysName, dbName string) element {
		return element{
			dbUser:       dbName,
			pattern:      regexp.MustCompile(sysName),
			substituteAt: strings.Index(dbName, `\1`),
		}
	}

	tcs := []struct {
		elt       element
		principal string
		expected  string
	}{
		{
			elt:       exactMatch("carlito", "carl"),
			principal: "carlito",
			expected:  "carl",
		},
		{
			elt:       exactMatch("carlito", "carl"),
			principal: "nope",
			expected:  "",
		},
		{
			elt:       regexMatch("^(.*)@cockroachlabs.com$", `\1`),
			principal: "carl@cockroachlabs.com",
			expected:  "carl",
		},
		{
			elt:       regexMatch("^(.*)@cockroachlabs.com$", `\11`),
			principal: "carl@cockroachlabs.com",
			expected:  "carl1",
		},
		{
			elt:       regexMatch("^(.*)@cockroachlabs.com$", `1\1`),
			principal: "carl@cockroachlabs.com",
			expected:  "1carl",
		},
		{
			elt:       regexMatch("^(.*)@cockroachlabs.com$", `\1`),
			principal: "carl@example.com",
			expected:  "",
		},
	}

	for idx, tc := range tcs {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			a := assert.New(t)
			a.Equal(tc.expected, tc.elt.substitute(tc.principal))
		})
	}
}

func TestIdentityMap(t *testing.T) {
	a := assert.New(t)
	data := `
# This is a comment
map-name system-username              database-username
foo      /^(.*)@cockroachlabs.com$    \1
foo      carl@cockroachlabs.com       also_carl   # Trailing comment
foo      carl@cockroachlabs.com       carl        # Duplicate behavior
`

	m, err := From(strings.NewReader(data))
	if !a.NoError(err) {
		return
	}
	t.Log(m.String())
	a.Len(m.data, 2)

	a.Nil(m.Map("missing", "carl"))

	if elts := m.data["map-name"]; a.Len(elts, 1) {
		a.Equal("database-username", elts[0].substitute("system-username"))
	}

	if elts, err := m.Map("foo", "carl@cockroachlabs.com"); a.NoError(err) && a.Len(elts, 2) {
		a.Equal("carl", elts[0].Normalized())
		a.Equal("also_carl", elts[1].Normalized())
	}

	a.Nil(m.Map("foo", "carl@example.com"))
}
