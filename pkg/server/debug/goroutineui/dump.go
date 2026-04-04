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

package goroutineui

import (
	"bytes"
	"io"
	"io/ioutil"
	"runtime"
	"sort"
	"strings"

	"github.com/maruel/panicparse/v2/stack"
)

// stacks is a wrapper for runtime.Stack that attempts to recover the data for all goroutines.
func stacks() []byte {
	// We don't know how big the traces are, so grow a few times if they don't fit. Start large, though.
	var trace []byte
	for n := 1 << 20; /* 1mb */ n <= (1 << 29); /* 512mb */ n *= 2 {
		trace = make([]byte, n)
		nbytes := runtime.Stack(trace, true /* all */)
		if nbytes < len(trace) {
			return trace[:nbytes]
		}
	}
	return trace
}

// A Dump wraps a goroutine dump with functionality to output through panicparse.
type Dump struct {
	agg *stack.Aggregated
	err error
}

// NewDump grabs a goroutine dump.
func NewDump() Dump {
	return newDumpFromBytes(stacks(), stack.DefaultOpts())
}

// newDumpFromBytes is like NewDump, but treats the supplied bytes as a goroutine
// dump. The function accepts the options to pass to panicparse/stack.ScanSnapshot.
func newDumpFromBytes(b []byte, opts *stack.Opts) Dump {
	s, _, err := stack.ScanSnapshot(bytes.NewBuffer(b), ioutil.Discard, opts)
	if err != io.EOF {
		return Dump{err: err}
	}
	return Dump{agg: s.Aggregate(stack.AnyValue)}
}

// SortCountDesc rearranges the goroutine buckets such that higher multiplicities
// appear earlier.
func (d Dump) SortCountDesc() {
	if d.err != nil {
		return
	}
	sort.Slice(d.agg.Buckets, func(i, j int) bool {
		a, b := d.agg.Buckets[i], d.agg.Buckets[j]
		return len(a.IDs) > len(b.IDs)
	})
}

// SortWaitDesc rearranges the goroutine buckets such that goroutines that have
// longer wait times appear earlier.
func (d Dump) SortWaitDesc() {
	if d.err != nil {
		return
	}
	sort.Slice(d.agg.Buckets, func(i, j int) bool {
		a, b := d.agg.Buckets[i], d.agg.Buckets[j]
		return a.SleepMax > b.SleepMax
	})
}

// HTML writes the rendered output of panicparse into the supplied Writer.
func (d Dump) HTML(w io.Writer) error {
	if d.err != nil {
		return d.err
	}
	return d.agg.ToHTML(w, "" /* footer */)
}

// HTMLString is like HTML, but returns a string. If an error occurs, its string
// representation is returned.
func (d Dump) HTMLString() string {
	var w strings.Builder
	if err := d.HTML(&w); err != nil {
		return err.Error()
	}
	return w.String()
}
