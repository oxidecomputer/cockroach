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

package spanconfigstore

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/protoutil"
)

type spanConfigKey string

// interner interns span config protos. It's a ref-counted data structure that
// maintains a single copy of every unique config that that's been added to
// it. Configs can also be removed, and if there are no more references to it,
// we no longer hold onto it in memory.
type interner struct {
	configToCanonical map[spanConfigKey]*roachpb.SpanConfig
	refCounts         map[*roachpb.SpanConfig]uint64
}

func newInterner() *interner {
	return &interner{
		configToCanonical: make(map[spanConfigKey]*roachpb.SpanConfig),
		refCounts:         make(map[*roachpb.SpanConfig]uint64),
	}
}

func (i *interner) add(ctx context.Context, conf roachpb.SpanConfig) *roachpb.SpanConfig {
	marshalled, err := protoutil.Marshal(&conf)
	if err != nil {
		log.Fatalf(ctx, "%v", err)
	}

	if canonical, found := i.configToCanonical[spanConfigKey(marshalled)]; found {
		i.refCounts[canonical]++
		return canonical
	}

	i.configToCanonical[spanConfigKey(marshalled)] = &conf
	i.refCounts[&conf] = 1
	return &conf
}

func (i *interner) remove(ctx context.Context, conf *roachpb.SpanConfig) {
	if _, found := i.refCounts[conf]; !found {
		return // nothing to do
	}

	i.refCounts[conf]--
	if i.refCounts[conf] != 0 {
		return // nothing to do
	}

	marshalled, err := protoutil.Marshal(conf)
	if err != nil {
		log.Fatalf(ctx, "%v", err)
	}

	delete(i.refCounts, conf)
	delete(i.configToCanonical, spanConfigKey(marshalled))
}

func (i *interner) copy() *interner {
	copiedInterner := newInterner()
	for k, v := range i.configToCanonical {
		copiedInterner.configToCanonical[k] = v
	}
	for k, v := range i.refCounts {
		copiedInterner.refCounts[k] = v
	}
	return copiedInterner
}
