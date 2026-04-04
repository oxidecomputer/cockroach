// Copyright 2017 The Cockroach Authors.
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

//go:build linux
// +build linux

package hlc

/*
#include <time.h>
*/
import "C"

import (
	"context"
	"os"

	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/errors"
)

// ClockSource contains the handle of the clock device as well as the
// clock id.
type ClockSource struct {
	// clockDevice is not used after the device is open but is here to prevent the GC
	// from closing the device and invalidating the clockDeviceID.
	clockDevice   *os.File
	clockDeviceID uintptr
}

// MakeClockSource creates a new ClockSource for the given device path.
func MakeClockSource(ctx context.Context, clockDevicePath string) (ClockSource, error) {
	var result ClockSource
	var err error
	result.clockDevice, err = os.Open(clockDevicePath)
	if err != nil {
		return result, errors.Wrapf(err, "cannot open %s", clockDevicePath)
	}

	clockDeviceFD := result.clockDevice.Fd()
	// For clarification of how the clock id is computed:
	// https://lore.kernel.org/patchwork/patch/868609/
	// https://github.com/torvalds/linux/blob/7e63420847ae5f1036e4f7c42f0b3282e73efbc2/tools/testing/selftests/ptp/testptp.c#L87
	clockID := (^clockDeviceFD << 3) | 3
	log.Infof(
		ctx,
		"opened clock device %s with fd %d, mod_fd %x",
		clockDevicePath,
		clockDeviceFD,
		clockID,
	)
	var ts C.struct_timespec
	_, err = C.clock_gettime(C.clockid_t(clockID), &ts)
	if err != nil {
		return result, errors.Wrap(err, "UseClockDevice: error calling clock_gettime")
	}
	result.clockDeviceID = clockID

	return result, nil
}

// UnixNano returns the clock device's physical nanosecond
// unix epoch timestamp as a convenience to create a HLC via
// c := hlc.NewClock(dev.UnixNano, ...).
func (p ClockSource) UnixNano() int64 {
	var ts C.struct_timespec
	_, err := C.clock_gettime(C.clockid_t(p.clockDeviceID), &ts)
	if err != nil {
		panic(err)
	}

	return int64(ts.tv_sec)*1e9 + int64(ts.tv_nsec)
}
