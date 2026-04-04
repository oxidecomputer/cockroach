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

// Package ptprovider encapsulates the concrete implementation of the
// protectedts.Provider.
package ptprovider

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/kv"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/protectedts"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/protectedts/ptcache"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/protectedts/ptreconcile"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/protectedts/ptstorage"
	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlutil"
	"github.com/cockroachdb/cockroach/pkg/util/metric"
	"github.com/cockroachdb/cockroach/pkg/util/stop"
	"github.com/cockroachdb/errors"
)

// Config configures the Provider.
type Config struct {
	Settings             *cluster.Settings
	DB                   *kv.DB
	Stores               *kvserver.Stores
	ReconcileStatusFuncs ptreconcile.StatusFuncs
	InternalExecutor     sqlutil.InternalExecutor
	Knobs                *protectedts.TestingKnobs
}

// Provider is the concrete implementation of protectedts.Provider interface.
type Provider struct {
	protectedts.Storage
	protectedts.Cache
	protectedts.Reconciler
	metric.Struct
}

// New creates a new protectedts.Provider.
func New(cfg Config) (protectedts.Provider, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}
	storage := ptstorage.New(cfg.Settings, cfg.InternalExecutor, cfg.Knobs)
	reconciler := ptreconcile.New(cfg.Settings, cfg.DB, storage, cfg.ReconcileStatusFuncs)
	cache := ptcache.New(ptcache.Config{
		DB:       cfg.DB,
		Storage:  storage,
		Settings: cfg.Settings,
	})

	return &Provider{
		Storage:    storage,
		Cache:      cache,
		Reconciler: reconciler,
		Struct:     reconciler.Metrics(),
	}, nil
}

func validateConfig(cfg Config) error {
	switch {
	case cfg.Settings == nil:
		return errors.Errorf("invalid nil Settings")
	case cfg.DB == nil:
		return errors.Errorf("invalid nil DB")
	case cfg.InternalExecutor == nil:
		return errors.Errorf("invalid nil InternalExecutor")
	default:
		return nil
	}
}

// Start implements the protectedts.Provider interface.
func (p *Provider) Start(ctx context.Context, stopper *stop.Stopper) error {
	if cache, ok := p.Cache.(*ptcache.Cache); ok {
		return cache.Start(ctx, stopper)
	}
	return nil
}

// Metrics implements the protectedts.Provider interface.
func (p *Provider) Metrics() metric.Struct {
	return p.Struct
}
