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

// Usage: go run ./listen_file_creation.go parent_folder_path file_name [timeout_duration]

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/fsnotify/fsnotify"
)

type result struct {
	finished bool
	err      error
}

const defaultTimeout = 30 * time.Second

func main() {
	if len(os.Args) < 2 {
		panic(errors.Wrap(
			fmt.Errorf("must provide the folder to watch and the file to listen to"),
			"fail to run fsnotify to listen to file creation"),
		)
	}

	var err error

	folderPath := os.Args[1]
	wantedFileName := os.Args[2]

	timeout := defaultTimeout

	var timeoutVal int
	if len(os.Args) > 3 {
		timeoutVal, err = strconv.Atoi(os.Args[3])
		if err != nil {
			panic(errors.Wrap(err, "timeout argument must be an integer"))
		}
	}

	timeout = time.Duration(timeoutVal) * time.Second

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(errors.Wrap(err, "cannot create new fsnotify file watcher"))
	}
	defer func() {
		if err := watcher.Close(); err != nil {
			panic(errors.Wrap(err, "error closing the file watcher in docker-fsnotify"))
		}
	}()

	done := make(chan result)

	go func() {
		for {
			if _, err := os.Stat(filepath.Join(folderPath, wantedFileName)); errors.Is(err, os.ErrNotExist) {
			} else {
				done <- result{finished: true, err: nil}
			}
			time.Sleep(time.Second * 1)
		}
	}()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fileName := event.Name[strings.LastIndex(event.Name, "/")+1:]
				if event.Op&fsnotify.Write == fsnotify.Write && fileName == wantedFileName {
					done <- result{finished: true, err: nil}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				done <- result{finished: false, err: err}
			}
		}
	}()

	err = watcher.Add(folderPath)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	select {
	case res := <-done:
		if res.finished && res.err == nil {
			fmt.Println("finished")
		} else {
			fmt.Printf("error in docker-fsnotify: %v", res.err)
		}

	case <-time.After(timeout):
		fmt.Printf("timeout for %s", timeout)
	}
}
