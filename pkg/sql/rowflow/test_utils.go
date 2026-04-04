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

package rowflow

import (
	"context"
	"sync"

	"github.com/cockroachdb/cockroach/pkg/sql/execinfra"
	"github.com/cockroachdb/cockroach/pkg/sql/execinfrapb"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
)

// MakeTestRouter creates a router to be used by tests.
func MakeTestRouter(
	ctx context.Context,
	flowCtx *execinfra.FlowCtx,
	spec *execinfrapb.OutputRouterSpec,
	streams []execinfra.RowReceiver,
	types []*types.T,
	wg *sync.WaitGroup,
) (execinfra.RowReceiver, error) {
	r, err := makeRouter(spec, streams)
	if err != nil {
		return nil, err
	}
	r.init(ctx, flowCtx, types)
	r.Start(ctx, wg, nil /* flowCtxCancel */)
	return r, nil
}
