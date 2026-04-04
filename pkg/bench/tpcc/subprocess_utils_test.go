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
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/testutils/skip"
	"github.com/cockroachdb/cockroach/pkg/util/envutil"
	"github.com/stretchr/testify/require"
)

type cmd struct {
	name string
	impl func(t *testing.T)
}

func (c *cmd) exec(env cmdEnv, args ...string) *exec.Cmd {
	cmd := exec.Command(os.Args[0],
		"--test.run=^TestInternal"+c.name+"$",
		"--test.v")
	if len(args) > 0 {
		cmd.Args = append(cmd.Args, "--")
		cmd.Args = append(cmd.Args, args...)
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=t", internalTestEnvVar))
	cmd.Env = append(cmd.Env, env.toStrings()...)
	return cmd
}

var isInternalTest = envutil.EnvOrDefaultBool(internalTestEnvVar, false)

func internalCommand(t *testing.T) {
	if !isInternalTest {
		skip.IgnoreLint(t)
	}
	f, ok := commands[strings.TrimPrefix(t.Name(), "TestInternal")]
	require.True(t, ok)
	f.impl(t)
}

func registerCmd(name string, impl func(t *testing.T)) *cmd {
	c := &cmd{
		name: name,
		impl: impl,
	}
	commands[name] = c
	return c
}

type envVar struct {
	k string
	v interface{}
}

type cmdEnv []envVar

func (ce cmdEnv) toStrings() (ret []string) {
	for _, v := range ce {
		ret = append(ret, fmt.Sprintf("%s=%v", v.k, v.v))
	}
	return ret
}
