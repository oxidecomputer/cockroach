// Copyright 2020 The Cockroach Authors.
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
	"os"
	"sort"
	"text/template"

	"github.com/cockroachdb/cockroach/pkg/cli/exit"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/spf13/cobra"
)

func init() {
	cmds = append(cmds, &cobra.Command{
		Use:   "logformats",
		Short: "Generate the markdown documentation for logging formats.",
		Args:  cobra.MaximumNArgs(1),
		Run:   runLogFormats,
	})
}

func runLogFormats(_ *cobra.Command, args []string) {
	if err := runLogFormatsInternal(args); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		exit.WithCode(exit.UnspecifiedError())
	}
}

func runLogFormatsInternal(args []string) error {
	// Compile the template.
	tmpl, err := template.New("format docs").Parse(fmtDocTemplate)
	if err != nil {
		return err
	}

	m := log.GetFormatterDocs()

	// Sort the names.
	fNames := make([]string, 0, len(m))
	for k := range m {
		fNames = append(fNames, k)
	}
	sort.Strings(fNames)

	// Retrieve the metadata into a format that the templating engine can understand.
	type info struct {
		Name string
		Doc  string
	}
	var infos []info
	for _, k := range fNames {
		infos = append(infos, info{Name: k, Doc: m[k]})
	}

	// Render the template.
	var src bytes.Buffer
	if err := tmpl.Execute(&src, struct {
		Formats []info
	}{infos}); err != nil {
		return err
	}

	// Write the output file.
	w := os.Stdout
	if len(args) > 0 {
		f, err := os.OpenFile(args[0], os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		defer func() { _ = f.Close() }()
		w = f
	}
	if _, err := w.Write(src.Bytes()); err != nil {
		return err
	}

	return nil
}

const fmtDocTemplate = `
The supported log output formats are documented below.

{{range .Formats}}
- [` + "`{{.Name}}`" + `](#format-{{.Name}})
{{end}}

{{range .Formats}}
## Format ` + "`{{.Name}}`" + `

{{.Doc}}
{{end}}
`
