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

package democluster

import (
	"time"

	"github.com/cockroachdb/cockroach/pkg/cli/clicfg"
	"github.com/cockroachdb/cockroach/pkg/workload"
)

// Context represents the input configuration and current state
// of a demo cluster.
type Context struct {
	// CliCtx links this demo context to a CLI configuration
	// environment.
	CliCtx *clicfg.Context

	// NumNodes is the requested number of nodes, and is also
	// modified when adding new nodes.
	NumNodes int

	// SQLPoolMemorySize is the size of the memory pool for each SQL
	// server.
	SQLPoolMemorySize int64

	// CacheSize is the size of the storage cache for each KV server.
	CacheSize int64

	// NoExampleDatabase prevents the auto-creation of a demo database
	// from a workload.
	NoExampleDatabase bool

	// RunWorkload indicates whether to run a workload in the background
	// after the demo cluster has been initialized.
	RunWorkload bool

	// WorkloadGenerator is the desired workload generator.
	WorkloadGenerator workload.Generator

	// WorkloadMaxQPS controls the amount of queries that can be run per
	// second.
	WorkloadMaxQPS int

	// Localities configures the list of localities available for use
	// by instantiated servers.
	Localities DemoLocalityList

	// GeoPartitionedReplicas requests that the executed workload
	// partition its data across localities. Requires an enterprise
	// license.
	GeoPartitionedReplicas bool

	// SimulateLatency requests that cross-region latencies be simulated
	// across region localities.
	SimulateLatency bool

	// DefaultKeySize is the default size of TLS private keys to use.
	DefaultKeySize int

	// DefaultCALifetime is the default lifetime of CA certs that are
	// generated for the transient cluster.
	DefaultCALifetime time.Duration

	// DefaultCertLifetime is the default lifetime of client certs that
	// are generated for the transient cluster.
	DefaultCertLifetime time.Duration

	// insecure requests that the server be started in "insecure mode".
	// NB: This is obsolete.
	Insecure bool

	// SQLPort is the first SQL port number to use when instantiating
	// servers. Use zero for auto-allocated random ports.
	SQLPort int

	// HTTPPort is the first HTTP port number to use when instantiating
	// servers. Use zero for auto-allocated random ports.
	HTTPPort int

	// ListeningURLFile can be set to a file which is written to after
	// the demo cluster has started, to contain a valid connection URL.
	ListeningURLFile string

	// Multitenant is true if we're starting the demo cluster in
	// multi-tenant mode.
	Multitenant bool

	// DefaultEnableRangefeeds is true if rangefeeds should start
	// out enabled.
	DefaultEnableRangefeeds bool
}

// IsInteractive returns true if the demo cluster configuration
// is for an interactive session. This exposes the field
// from clicfg.Context if available.
func (demoCtx *Context) IsInteractive() bool {
	return demoCtx.CliCtx != nil && demoCtx.CliCtx.IsInteractive
}
