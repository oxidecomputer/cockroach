// Copyright 2015 The Cockroach Authors.
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

package log

import (
	"context"
	"fmt"
	"strings"
)

// FormatWithContextTags formats the string and prepends the context
// tags.
//
// Redaction markers are *not* inserted. The resulting
// string is generally unsafe for reporting.
func FormatWithContextTags(ctx context.Context, format string, args ...interface{}) string {
	var buf strings.Builder
	formatTags(ctx, true /* brackets */, &buf)
	formatArgs(&buf, format, args...)
	return buf.String()
}

func formatOnlyArgs(format string, args ...interface{}) string {
	var buf strings.Builder
	formatArgs(&buf, format, args...)
	return buf.String()
}

func formatArgs(buf *strings.Builder, format string, args ...interface{}) {
	if len(args) == 0 {
		buf.WriteString(format)
	} else if len(format) == 0 {
		fmt.Fprint(buf, args...)
	} else {
		fmt.Fprintf(buf, format, args...)
	}
}
