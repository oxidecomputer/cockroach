// Copyright 2018 The Cockroach Authors.
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

package ssh

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestProgress(t *testing.T) {
	output, err := ioutil.TempFile("", "example*")
	if err != nil {
		t.Fatal(err)
	}
	defer output.Close()
	defer func() {
		if err := os.Remove(output.Name()); err != nil {
			t.Fatal(err)
		}
	}()

	b := make([]byte, 10)
	var percent float64
	writer := &ProgressWriter{
		Writer: output,
		Done:   0,
		Total:  50,
		Progress: func(currentProgress float64) {
			percent = currentProgress
		},
	}
	for i := 0; i < 4; i++ {
		if _, err := writer.Write(b); err != nil {
			t.Fatal(err)
		}
	}
	if percent != 0.8 {
		t.Errorf("expected progress of 80%% but got %.2f", percent*100)
	}
	if _, err := writer.Write(b); err != nil {
		t.Fatal(err)
	}
	if percent != 1.0 {
		t.Errorf("expected progress of 100%% but got %.2f", percent*100)
	}
}
