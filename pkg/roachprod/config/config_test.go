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

package config

import "testing"

func TestIsLocalClusterName(t *testing.T) {
	yes := []string{
		"local",
		"local-1",
		"local-foo",
		"local-foo-bar-123-aZy",
	}
	no := []string{
		"loca",
		"locall",
		"local1",
		"local-",
		"local-foo?",
		"local-foo/",
	}

	for _, s := range yes {
		if !IsLocalClusterName(s) {
			t.Errorf("expected '%s' to be a valid local cluster name", s)
		}
	}

	for _, s := range no {
		if IsLocalClusterName(s) {
			t.Errorf("expected '%s' to not be a valid local cluster name", s)
		}
	}
}
