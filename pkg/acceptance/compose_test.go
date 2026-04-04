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

package acceptance

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/acceptance/cluster"
	"github.com/cockroachdb/cockroach/pkg/build/bazel"
)

const composeDir = "compose"

func TestComposeGSS(t *testing.T) {
	testCompose(t, filepath.Join("gss", "docker-compose.yml"), "psql")
}

func TestComposeGSSPython(t *testing.T) {
	testCompose(t, filepath.Join("gss", "docker-compose-python.yml"), "python")
}

func TestComposeFlyway(t *testing.T) {
	testCompose(t, filepath.Join("flyway", "docker-compose.yml"), "flyway")
}

func testCompose(t *testing.T, path string, exitCodeFrom string) {
	if bazel.BuiltWithBazel() {
		// Copy runfiles symlink content to a temporary directory to avoid broken symlinks in docker.
		tmpComposeDir, err := ioutil.TempDir("", "")
		if err != nil {
			t.Fatal(err.Error())
		}
		err = copyRunfiles(composeDir, tmpComposeDir)
		if err != nil {
			t.Fatal(err.Error())
		}
		defer func() {
			_ = os.RemoveAll(tmpComposeDir)
		}()
		path = filepath.Join(tmpComposeDir, path)
		// If running under Bazel, export 2 environment variables that will be interpolated in docker-compose.yml files.
		cockroachBinary, err := filepath.Abs(*cluster.CockroachBinary)
		if err != nil {
			t.Fatal(err.Error())
		}
		err = os.Setenv("COCKROACH_BINARY", cockroachBinary)
		if err != nil {
			t.Fatal(err.Error())
		}
		err = os.Setenv("CERTS_DIR", cluster.AbsCertsDir())
		if err != nil {
			t.Fatal(err.Error())
		}
	} else {
		path = filepath.Join(composeDir, path)
	}
	uid := os.Getuid()
	err := os.Setenv("UID", strconv.Itoa(uid))
	if err != nil {
		t.Fatal(err.Error())
	}
	gid := os.Getgid()
	err = os.Setenv("GID", strconv.Itoa(gid))
	if err != nil {
		t.Fatal(err.Error())
	}
	cmd := exec.Command(
		"docker-compose",
		"--no-ansi",
		"-f", path,
		"up",
		"--force-recreate",
		"--build",
		"--exit-code-from", exitCodeFrom,
	)
	var buf bytes.Buffer
	if testing.Verbose() {
		cmd.Stdout = io.MultiWriter(&buf, os.Stdout)
		cmd.Stderr = io.MultiWriter(&buf, os.Stderr)
	} else {
		cmd.Stdout = &buf
		cmd.Stderr = &buf
	}
	if err := cmd.Run(); err != nil {
		t.Log(buf.String())
		t.Fatal(err)
	}
}
