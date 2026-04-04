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

package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"

	bazelutil "github.com/cockroachdb/cockroach/pkg/build/util"
	"github.com/spf13/cobra"
)

func makeMergeTestXMLsCmd(runE func(cmd *cobra.Command, args []string) error) *cobra.Command {
	mergeTestXMLsCommand := &cobra.Command{
		Use:   "merge-test-xmls XML1 [XML2...]",
		Short: "Merge the given test XML's (utility command)",
		Long:  "Merge the given test XML's (utility command)",
		Args:  cobra.MinimumNArgs(1),
		RunE:  runE,
	}
	mergeTestXMLsCommand.Hidden = true
	return mergeTestXMLsCommand
}

func (d *dev) mergeTestXMLs(cmd *cobra.Command, xmls []string) error {
	var suites []bazelutil.TestSuites
	for _, file := range xmls {
		suitesToAdd := bazelutil.TestSuites{}
		input, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		err = xml.Unmarshal(input, &suitesToAdd)
		if err != nil {
			return err
		}
		suites = append(suites, suitesToAdd)
	}
	return bazelutil.MergeTestXMLs(suites, os.Stdout)
}
