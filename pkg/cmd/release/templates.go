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
	"bytes"
	"fmt"
	htmltemplate "html/template"
	"text/template"
)

// templateToText is helper function to execute a template using the text/template package
func templateToText(templateText string, args interface{}) (string, error) {
	templ, err := template.New("").Parse(templateText)
	if err != nil {
		return "", fmt.Errorf("cannot parse template: %w", err)
	}

	var buf bytes.Buffer
	err = templ.Execute(&buf, args)
	if err != nil {
		return "", fmt.Errorf("cannot execute template: %w", err)
	}
	return buf.String(), nil
}

// templateToHTML is helper function to execute a template using the html/template package
func templateToHTML(templateText string, args interface{}) (string, error) {
	templ, err := htmltemplate.New("").Parse(templateText)
	if err != nil {
		return "", fmt.Errorf("cannot parse template: %w", err)
	}

	var buf bytes.Buffer
	err = templ.Execute(&buf, args)
	if err != nil {
		return "", fmt.Errorf("cannot execute template: %w", err)
	}
	return buf.String(), nil
}
