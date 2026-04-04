// Copyright 2022 The Cockroach Authors.
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

package pprofutil

import (
	"context"
	"runtime/pprof"

	"github.com/cockroachdb/logtags"
)

// SetProfilerLabels returns a context wrapped with the provided pprof labels
// provided in alternating key-value format. The returned closure should be
// defer'ed to restore the original labels from the initial context.
//
// This method allocates and isn't suitable for hot code paths.
func SetProfilerLabels(ctx context.Context, labels ...string) (_ context.Context, undo func()) {
	origCtx := ctx
	ctx = pprof.WithLabels(ctx, pprof.Labels(labels...))
	pprof.SetGoroutineLabels(ctx)
	return ctx, func() {
		pprof.SetGoroutineLabels(origCtx)
	}
}

// SetProfilerLabelsFromCtxTags is like SetProfilerLabels, but sources the labels from
// the logtags.Buffer in the context, if any. If there is no buffer or it has no tags,
// the goroutine labels are not updated.
func SetProfilerLabelsFromCtxTags(ctx context.Context) (_ context.Context, undo func()) {
	tags := logtags.FromContext(ctx)
	if tags == nil || len(tags.Get()) == 0 {
		return ctx, func() {}
	}
	var labels []string
	for _, tag := range tags.Get() {
		labels = append(labels, tag.Key(), tag.ValueStr())
	}
	return SetProfilerLabels(ctx, labels...)
}
