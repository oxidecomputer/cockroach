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

package tpcc

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/base"
)

type option interface {
	fmt.Stringer
	apply(*benchmarkConfig)
}

type options []option

func (o options) String() string {
	var buf strings.Builder
	for i, opt := range o {
		if i > 0 {
			buf.WriteString(";")
		}
		buf.WriteString(opt.String())
	}
	return buf.String()
}

func (o options) apply(cfg *benchmarkConfig) {
	for _, opt := range o {
		opt.apply(cfg)
	}
}

type benchmarkConfig struct {
	workloadFlags []string
	argsGenerator serverArgs
	setupStmts    []string
}

type workloadFlagOption struct{ name, value string }

func (w workloadFlagOption) apply(cfg *benchmarkConfig) {
	cfg.workloadFlags = append(cfg.workloadFlags, "--"+w.name, w.value)
}

func (w workloadFlagOption) String() string {
	return fmt.Sprintf("%s=%s", w.name, w.value)
}

func workloadFlag(name, value string) option {
	return workloadFlagOption{name: name, value: value}
}

type serverArgs func(b testing.TB) (_ base.TestServerArgs, cleanup func())

func (s serverArgs) apply(cfg *benchmarkConfig) {
	cfg.argsGenerator = s
}

func (s serverArgs) String() string { return "generator" }

type setupStmtOption string

func (s setupStmtOption) apply(cfg *benchmarkConfig) {
	cfg.setupStmts = append(cfg.setupStmts, string(s))
}

func setupStmt(stmt string) option {
	return setupStmtOption(stmt)
}

func (s setupStmtOption) String() string { return string(s) }
