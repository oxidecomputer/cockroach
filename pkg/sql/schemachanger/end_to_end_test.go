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

package schemachanger_test

import (
	"flag"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/sctest"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/skip"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

// Used for saving corpus information in TestGenerateCorpus
var corpusPath string

func init() {
	flag.StringVar(&corpusPath, "declarative-corpus", "", "Path to the corpus file")
}

func TestSchemaChangerSideEffects(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.EndToEndSideEffects(t, testutils.TestDataPath(t), sctest.SingleNodeCluster)
}

func TestRollback(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	sctest.Rollback(t, testutils.TestDataPath(t), sctest.SingleNodeCluster)
}

func TestPause(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	sctest.Pause(t, testutils.TestDataPath(t), sctest.SingleNodeCluster)
}

// TestGenerateCorpus generates a corpus based on the end to end test files.
func TestGenerateCorpus(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	if corpusPath == "" {
		skip.IgnoreLintf(t, "requires decalarative-corpus path parameter")
	}
	sctest.GenerateSchemaChangeCorpus(t, testutils.TestDataPath(t), corpusPath, sctest.SingleNodeCluster)
}
