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

package sctestdeps

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scexec"
)

type backfillTrackerDeps interface {
	Catalog() scexec.Catalog
	Codec() keys.SQLCodec
}

type testBackfillTracker struct {
	deps backfillTrackerDeps
}

// BackfillProgressTracker implements the scexec.Dependencies interface.
func (s *TestState) BackfillProgressTracker() scexec.BackfillTracker {
	return s.backfillTracker
}

var _ scexec.BackfillTracker = (*testBackfillTracker)(nil)

func (s *testBackfillTracker) GetBackfillProgress(
	ctx context.Context, b scexec.Backfill,
) (scexec.BackfillProgress, error) {
	return scexec.BackfillProgress{Backfill: b}, nil
}

func (s *testBackfillTracker) SetBackfillProgress(
	ctx context.Context, progress scexec.BackfillProgress,
) error {
	return nil
}

func (s *testBackfillTracker) FlushCheckpoint(ctx context.Context) error {
	return nil
}

func (s *testBackfillTracker) FlushFractionCompleted(ctx context.Context) error {
	return nil
}

// StartPeriodicFlush implements the scexec.PeriodicProgressFlusher interface.
func (s *TestState) StartPeriodicFlush(
	ctx context.Context,
) (close func(context.Context) error, _ error) {
	return func(ctx context.Context) error { return nil }, nil
}

type testBackfiller struct {
	s *TestState
}

var _ scexec.Backfiller = (*testBackfiller)(nil)

// BackfillIndex implements the scexec.Backfiller interface.
func (s *testBackfiller) BackfillIndex(
	_ context.Context,
	progress scexec.BackfillProgress,
	_ scexec.BackfillProgressWriter,
	tbl catalog.TableDescriptor,
) error {
	s.s.LogSideEffectf(
		"backfill indexes %v from index #%d in table #%d",
		progress.DestIndexIDs, progress.SourceIndexID, tbl.GetID(),
	)
	return nil
}

// MaybePrepareDestIndexesForBackfill implements the scexec.Backfiller interface.
func (s *testBackfiller) MaybePrepareDestIndexesForBackfill(
	ctx context.Context, progress scexec.BackfillProgress, descriptor catalog.TableDescriptor,
) (scexec.BackfillProgress, error) {
	return progress, nil
}

var _ scexec.IndexSpanSplitter = (*indexSpanSplitter)(nil)

type indexSpanSplitter struct{}

// MaybeSplitIndexSpans implements the scexec.IndexSpanSplitter interface.
func (s *indexSpanSplitter) MaybeSplitIndexSpans(
	_ context.Context, _ catalog.TableDescriptor, _ catalog.Index,
) error {
	return nil
}
