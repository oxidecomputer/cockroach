// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package diagnostics

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/server/diagnostics/diagnosticspb"
	"github.com/cockroachdb/cockroach/pkg/util/cloudinfo"
	"github.com/cockroachdb/cockroach/pkg/util/syncutil"
	"github.com/cockroachdb/cockroach/pkg/util/system"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

var populateMutex syncutil.Mutex

// populateHardwareInfo populates OS, CPU, memory, etc. information about the
// environment in which CRDB is running.
func populateHardwareInfo(ctx context.Context, e *diagnosticspb.Environment) {
	// The shirou/gopsutil/host library is not multi-thread safe. As one
	// example, it lazily initializes a global map the first time the
	// Virtualization function is called, but takes no lock while doing so.
	// Work around this limitation by taking our own lock.
	populateMutex.Lock()
	defer populateMutex.Unlock()

	if platform, family, version, err := host.PlatformInformation(); err == nil {
		e.Os.Family = family
		e.Os.Platform = platform
		e.Os.Version = version
	}

	if virt, role, err := host.Virtualization(); err == nil && role == "guest" {
		e.Hardware.Virtualization = virt
	}

	if m, err := mem.VirtualMemory(); err == nil {
		e.Hardware.Mem.Available = m.Available
		e.Hardware.Mem.Total = m.Total
	}

	e.Hardware.Cpu.Numcpu = int32(system.NumCPU())
	if cpus, err := cpu.InfoWithContext(ctx); err == nil && len(cpus) > 0 {
		e.Hardware.Cpu.Sockets = int32(len(cpus))
		c := cpus[0]
		e.Hardware.Cpu.Cores = c.Cores
		e.Hardware.Cpu.Model = c.ModelName
		e.Hardware.Cpu.Mhz = float32(c.Mhz)
		e.Hardware.Cpu.Features = c.Flags
	}

	if l, err := load.AvgWithContext(ctx); err == nil {
		e.Hardware.Loadavg15 = float32(l.Load15)
	}

	e.Hardware.Provider, e.Hardware.InstanceClass = cloudinfo.GetInstanceClass(ctx)
	e.Topology.Provider, e.Topology.Region = cloudinfo.GetInstanceRegion(ctx)
}
