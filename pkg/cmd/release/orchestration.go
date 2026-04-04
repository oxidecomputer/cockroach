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

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	versionFlag      = "version"
	templatesDirFlag = "template-dir"
	outputDirFlag    = "output-dir"
)

var orchestrationFlags = struct {
	version      string
	templatesDir string
	outputDir    string
}{}

var setOrchestrationVersionCmd = &cobra.Command{
	Use:   "set-orchestration-version",
	Short: "Set orchestration version",
	Long:  "Updates orchestration version under the ./cloud/kubernetes directory",
	RunE:  setOrchestrationVersion,
}

func init() {
	setOrchestrationVersionCmd.Flags().StringVar(&orchestrationFlags.version, versionFlag, "", "cockroachdb version")
	setOrchestrationVersionCmd.Flags().StringVar(&orchestrationFlags.templatesDir, templatesDirFlag, "",
		"orchestration templates directory")
	setOrchestrationVersionCmd.Flags().StringVar(&orchestrationFlags.outputDir, outputDirFlag, "",
		"orchestration directory")
	_ = setOrchestrationVersionCmd.MarkFlagRequired(versionFlag)
	_ = setOrchestrationVersionCmd.MarkFlagRequired(templatesDirFlag)
	_ = setOrchestrationVersionCmd.MarkFlagRequired(outputDirFlag)
}

func setOrchestrationVersion(_ *cobra.Command, _ []string) error {
	dirInfo, err := os.Stat(orchestrationFlags.templatesDir)
	if err != nil {
		return fmt.Errorf("cannot stat templates directory: %w", err)
	}
	if !dirInfo.IsDir() {
		return fmt.Errorf("%s is not a directory", orchestrationFlags.templatesDir)
	}
	return filepath.Walk(orchestrationFlags.templatesDir, func(filePath string, fileInfo os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		// Skip directories
		if !fileInfo.Mode().IsRegular() {
			return nil
		}
		// calculate file directory relative to the given root directory.
		dir := path.Dir(filePath)
		relDir, err := filepath.Rel(orchestrationFlags.templatesDir, dir)
		if err != nil {
			return err
		}
		destDir := filepath.Join(orchestrationFlags.outputDir, relDir)
		destFile := filepath.Join(destDir, fileInfo.Name())
		if err := os.MkdirAll(destDir, 0755); err != nil && !errors.Is(err, os.ErrExist) {
			return err
		}
		contents, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}
		// Go templates cannot be used here, because some files are templates already.
		generatedContents := strings.ReplaceAll(string(contents), "@VERSION@", orchestrationFlags.version)
		if strings.HasSuffix(destFile, ".yaml") {
			generatedContents = fmt.Sprintf("# Generated file, DO NOT EDIT. Source: %s\n", filePath) + generatedContents
		}
		err = ioutil.WriteFile(destFile, []byte(generatedContents), fileInfo.Mode().Perm())
		if err != nil {
			return err
		}
		return nil
	})
}
