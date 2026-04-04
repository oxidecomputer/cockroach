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
package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseArgs(t *testing.T) {
	_, err := parseArgs([]string{}, -1)
	assert.True(t, errors.Is(err, errUsage))
	_, err = parseArgs([]string{"build"}, -1)
	assert.True(t, errors.Is(err, errUsage))
	_, err = parseArgs([]string{"typo", "target"}, -1)
	assert.NotNil(t, err)
	args, err := parseArgs([]string{"build", "target"}, -1)
	assert.Nil(t, err)
	assert.Equal(t, parsedArgs{
		subcmd:     "build",
		targets:    []string{"target"},
		additional: []string{}}, *args)
	args, err = parseArgs([]string{"test", "target1", "target2"}, -1)
	assert.Nil(t, err)
	assert.Equal(t, parsedArgs{
		subcmd:     "test",
		targets:    []string{"target1", "target2"},
		additional: []string{}}, *args)
	// Make sure additional arguments are captured correctly.
	args, err = parseArgs([]string{"test", "target1", "target2", "--verbose_failures"}, 3)
	assert.Nil(t, err)
	assert.Equal(t, parsedArgs{
		subcmd:     "test",
		targets:    []string{"target1", "target2"},
		additional: []string{"--verbose_failures"},
	}, *args)
}
