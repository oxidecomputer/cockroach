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

package corpus_test

import (
	"context"
	"flag"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/jobs/jobspb"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/corpus"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scop"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scplan"
	"github.com/cockroachdb/cockroach/pkg/testutils/skip"
	"github.com/stretchr/testify/require"
)

// Used for reading corpus information in TestValidateCorpuses
var corpusPath string

func init() {
	flag.StringVar(&corpusPath, "declarative-corpus", "", "path to the corpus file")
}

// TestValidateCorpuses validates that any generated corpus file on disk, a
// path needs to be specified.
func TestValidateCorpuses(t *testing.T) {
	if corpusPath == "" {
		skip.IgnoreLintf(t, "requires declarative-corpus path parameter")
	}
	reader, err := corpus.NewCorpusReader(corpusPath)
	require.NoError(t, err)
	require.NoError(t, reader.ReadCorpus())
	for corpusIdx := 0; corpusIdx < reader.GetNumEntries(); corpusIdx++ {
		jobID := jobspb.InvalidJobID
		name, state := reader.GetCorpus(corpusIdx)
		t.Run(name, func(t *testing.T) {
			_, err := scplan.MakePlan(context.Background(), *state, scplan.Params{
				ExecutionPhase: scop.LatestPhase,
				InRollback:     state.InRollback,
				SchemaChangerJobIDSupplier: func() jobspb.JobID {
					jobID++
					return jobID
				}})
			require.NoError(t, err)
		})
	}
}
