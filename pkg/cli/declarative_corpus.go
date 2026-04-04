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

package cli

import (
	"context"
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/cli/clierrorplus"
	"github.com/cockroachdb/cockroach/pkg/jobs/jobspb"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/corpus"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scop"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scplan"
	"github.com/spf13/cobra"
)

var declarativeValidateCorpus = &cobra.Command{
	Use:   "declarative-corpus-validate <filename>",
	Short: "validates a corpus file for the declarative schema changer",
	Long: `
Validates very single declarative schema changer state can be planned against in
a given corpus file.
`,
	Args: cobra.ExactArgs(1),
	RunE: clierrorplus.MaybeDecorateError(
		func(cmd *cobra.Command, args []string) (resErr error) {
			cr, err := corpus.NewCorpusReaderWithPath(args[0])
			if err != nil {
				panic(err)
			}
			err = cr.ReadCorpus()
			if err != nil {
				panic(err)
			}
			for idx := 0; idx < cr.GetNumEntries(); idx++ {
				name, state := cr.GetCorpus(idx)
				jobID := jobspb.JobID(0)
				params := scplan.Params{
					InRollback:     state.InRollback,
					ExecutionPhase: scop.LatestPhase,
					SchemaChangerJobIDSupplier: func() jobspb.JobID {
						jobID++
						return jobID
					},
				}
				_, err := scplan.MakePlan(context.Background(), *state, params)
				if err != nil {
					fmt.Printf("failed to validate %s with error %v\n", name, err)
				} else {
					fmt.Printf("validated %s\n", name)
				}
			}
			return nil
		}),
}
