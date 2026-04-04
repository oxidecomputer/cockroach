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

package server

import (
	"context"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/storage"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/require"
)

func TestStickyEngines(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	ctx := context.Background()
	attrs := roachpb.Attributes{}
	cacheSize := int64(1 << 20)   /* 1 MiB */
	storeSize := int64(512 << 20) /* 512 MiB */

	registry := NewStickyInMemEnginesRegistry()

	cfg1 := MakeConfig(ctx, cluster.MakeTestingClusterSettings())
	cfg1.CacheSize = cacheSize
	spec1 := base.StoreSpec{
		StickyInMemoryEngineID: "engine1",
		Attributes:             attrs,
		Size:                   base.SizeSpec{InBytes: storeSize},
	}
	engine1, err := registry.GetOrCreateStickyInMemEngine(ctx, &cfg1, spec1)
	require.NoError(t, err)
	require.False(t, engine1.Closed())

	cfg2 := MakeConfig(ctx, cluster.MakeTestingClusterSettings())
	cfg2.CacheSize = cacheSize
	spec2 := base.StoreSpec{
		StickyInMemoryEngineID: "engine2",
		Attributes:             attrs,
		Size:                   base.SizeSpec{InBytes: storeSize},
	}
	engine2, err := registry.GetOrCreateStickyInMemEngine(ctx, &cfg2, spec2)
	require.NoError(t, err)
	require.False(t, engine2.Closed())

	// Regetting the engine whilst it is not closed will fail.
	_, err = registry.GetOrCreateStickyInMemEngine(ctx, &cfg1, spec1)
	require.EqualError(t, err, "sticky engine engine1 has not been closed")

	// Close the engine, which allows it to be refetched.
	engine1.Close()
	require.True(t, engine1.Closed())
	require.False(t, engine1.(*stickyInMemEngine).Engine.Closed())

	// Refetching the engine should give back the same engine.
	engine1Refetched, err := registry.GetOrCreateStickyInMemEngine(ctx, &cfg1, spec1)
	require.NoError(t, err)
	require.Equal(t, engine1, engine1Refetched)
	require.False(t, engine1.Closed())

	// Cleaning up everything asserts everything is closed.
	registry.CloseAllStickyInMemEngines()
	for _, engine := range []storage.Engine{engine1, engine2} {
		require.True(t, engine.Closed())
		require.True(t, engine.(*stickyInMemEngine).Engine.Closed())
	}
}

func TestStickyEnginesReplaceEngines(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	ctx := context.Background()
	attrs := roachpb.Attributes{}
	cacheSize := int64(1 << 20)   /* 1 MiB */
	storeSize := int64(512 << 20) /* 512 MiB */

	registry := NewStickyInMemEnginesRegistry(ReplaceEngines)

	cfg1 := MakeConfig(ctx, cluster.MakeTestingClusterSettings())
	cfg1.CacheSize = cacheSize
	spec1 := base.StoreSpec{
		StickyInMemoryEngineID: "engine1",
		Attributes:             attrs,
		Size:                   base.SizeSpec{InBytes: storeSize},
	}
	engine1, err := registry.GetOrCreateStickyInMemEngine(ctx, &cfg1, spec1)
	require.NoError(t, err)
	fs1, err := registry.GetUnderlyingFS(spec1)
	require.NoError(t, err)
	require.False(t, engine1.Closed())
	engine1.Close()

	// Refetching the engine should give back a different engine with the same
	// underlying fs.
	engine1Refetched, err := registry.GetOrCreateStickyInMemEngine(ctx, &cfg1, spec1)
	require.NoError(t, err)
	fs1Refetched, err := registry.GetUnderlyingFS(spec1)
	require.NoError(t, err)
	require.NotEqual(t, engine1, engine1Refetched)
	require.Equal(t, fs1, fs1Refetched)
	require.True(t, engine1.Closed())
	require.False(t, engine1Refetched.Closed())

	registry.CloseAllStickyInMemEngines()
	for _, engine := range []storage.Engine{engine1, engine1Refetched} {
		require.True(t, engine.Closed())
		require.True(t, engine.(*stickyInMemEngine).Engine.Closed())
	}
}
