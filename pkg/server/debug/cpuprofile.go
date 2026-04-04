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

package debug

import (
	"net/http"
	"net/http/pprof"
	"strconv"

	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
)

// CPUProfileOptions contains options for generating a CPU profile.
type CPUProfileOptions struct {
	// Number of seconds to profile for.
	Seconds int32
	// Whether to enable pprof labels while the profile is taken.
	WithLabels bool
}

// Type returns the CPUProfileType corresponding to the options.
func (opts CPUProfileOptions) Type() cluster.CPUProfileType {
	typ := cluster.CPUProfileDefault
	if opts.WithLabels {
		typ = cluster.CPUProfileWithLabels
	}
	return typ
}

// CPUProfileOptionsFromRequest parses the `seconds` and `labels` fragments
// from the URL and populates CPUProfileOptions from it.
//
// For convenience, `labels` defaults to true, that is, `?labels=false`
// must be specified to disable them. `seconds` defaults to the pprof
// default of 30s.
func CPUProfileOptionsFromRequest(r *http.Request) CPUProfileOptions {
	seconds, err := strconv.ParseInt(r.FormValue("seconds"), 10, 32)
	if err != nil || seconds <= 0 {
		seconds = 30
	}
	// NB: default to using labels unless it's specifically set to false.
	withLabels := r.FormValue("labels") != "false"
	return CPUProfileOptions{
		Seconds:    int32(seconds),
		WithLabels: withLabels,
	}
}

// CPUProfileDo invokes the closure while enabling (and disabling) the supplied
// CPUProfileMode. Errors if the profiling mode could not be set or if do()
// returns an error.
func CPUProfileDo(st *cluster.Settings, typ cluster.CPUProfileType, do func() error) error {
	if err := st.SetCPUProfiling(typ); err != nil {
		return err
	}
	defer func() { _ = st.SetCPUProfiling(cluster.CPUProfileNone) }()
	return do()
}

// CPUProfileHandler is replacement for `pprof.Profile` that supports additional
// options.
func CPUProfileHandler(st *cluster.Settings, w http.ResponseWriter, r *http.Request) {
	opts := CPUProfileOptionsFromRequest(r)
	if err := CPUProfileDo(st, opts.Type(), func() error {
		pprof.Profile(w, r)
		return nil
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
