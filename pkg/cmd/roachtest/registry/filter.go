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

package registry

import (
	"regexp"
	"strings"
)

// TestFilter holds the name and tag filters for filtering tests.
// See NewTestFilter.
type TestFilter struct {
	Name *regexp.Regexp
	Tag  *regexp.Regexp
	// RawTag is the string representation of the regexps in tag.
	RawTag []string
}

// NewTestFilter initializes a new filter. The strings are interpreted
// as regular expressions. As a special case, a `tag:` prefix implies
// that the remainder of the string filters tests by tag, and not by
// name.
func NewTestFilter(filter []string) *TestFilter {
	var name []string
	var tag []string
	var rawTag []string
	for _, v := range filter {
		if strings.HasPrefix(v, "tag:") {
			tag = append(tag, strings.TrimPrefix(v, "tag:"))
			rawTag = append(rawTag, v)
		} else {
			name = append(name, v)
		}
	}

	if len(tag) == 0 {
		tag = []string{DefaultTag}
		rawTag = []string{"tag:" + DefaultTag}
	}

	makeRE := func(strs []string) *regexp.Regexp {
		switch len(strs) {
		case 0:
			return regexp.MustCompile(`.`)
		case 1:
			return regexp.MustCompile(strs[0])
		default:
			for i := range strs {
				strs[i] = "(" + strs[i] + ")"
			}
			return regexp.MustCompile(strings.Join(strs, "|"))
		}
	}

	return &TestFilter{
		Name:   makeRE(name),
		Tag:    makeRE(tag),
		RawTag: rawTag,
	}
}
