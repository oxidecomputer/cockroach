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

package humanizeutil

import (
	"github.com/cockroachdb/redact"
	"github.com/dustin/go-humanize"
)

// Count formats a unitless integer value like a row count. It uses separating
// commas for large values (e.g. "1,000,000").
func Count(val uint64) redact.SafeString {
	return redact.SafeString(humanize.Comma(int64(val)))
}
