// Copyright 2019 The Cockroach Authors.
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

package rangefeed

import (
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/util/interval"
)

// Filter informs the producer of logical operations of the information that a
// rangefeed Processor is interested in, given its current set of registrations.
// It can be used to avoid performing extra work to provide the Processor with
// information which will be ignored.
type Filter struct {
	needPrevVals interval.RangeGroup
	needVals     interval.RangeGroup
}

func newFilterFromRegistry(reg *registry) *Filter {
	f := &Filter{
		needPrevVals: interval.NewRangeList(),
		needVals:     interval.NewRangeList(),
	}
	reg.tree.Do(func(i interval.Interface) (done bool) {
		r := i.(*registration)
		if r.withDiff {
			f.needPrevVals.Add(r.Range())
		}
		f.needVals.Add(r.Range())
		return false
	})
	return f
}

// NeedPrevVal returns whether the Processor requires MVCCWriteValueOp and
// MVCCCommitIntentOp operations over the specified key span to contain
// populated PrevValue fields.
func (r *Filter) NeedPrevVal(s roachpb.Span) bool {
	return r.needPrevVals.Overlaps(s.AsRange())
}

// NeedVal returns whether the Processor requires MVCCWriteValueOp and
// MVCCCommitIntentOp operations over the specified key span to contain
// populated Value fields.
func (r *Filter) NeedVal(s roachpb.Span) bool {
	return r.needVals.Overlaps(s.AsRange())
}
