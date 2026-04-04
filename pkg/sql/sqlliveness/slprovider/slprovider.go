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

// Package slprovider exposes an implementation of the sqlliveness.Provider
// interface.
package slprovider

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/kv"
	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlliveness"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlliveness/slinstance"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlliveness/slstorage"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/metric"
	"github.com/cockroachdb/cockroach/pkg/util/stop"
)

// New constructs a new Provider.
func New(
	ambientCtx log.AmbientContext,
	stopper *stop.Stopper,
	clock *hlc.Clock,
	db *kv.DB,
	codec keys.SQLCodec,
	settings *cluster.Settings,
	testingKnobs *sqlliveness.TestingKnobs,
) sqlliveness.Provider {
	storage := slstorage.NewStorage(ambientCtx, stopper, clock, db, codec, settings)
	instance := slinstance.NewSQLInstance(stopper, clock, storage, settings, testingKnobs)
	return &provider{
		Storage:  storage,
		Instance: instance,
	}
}

func (p *provider) Start(ctx context.Context) {
	p.Storage.Start(ctx)
	p.Instance.Start(ctx)
}

func (p *provider) Metrics() metric.Struct {
	return p.Storage.Metrics()
}

type provider struct {
	*slstorage.Storage
	*slinstance.Instance
}
