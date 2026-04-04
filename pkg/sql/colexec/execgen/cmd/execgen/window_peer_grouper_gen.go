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
	"io"
	"strings"
	"text/template"
)

type windowPeerGrouperTmplInfo struct {
	AllPeers     bool
	HasPartition bool
	String       string
}

const windowPeerGrouperOpsTmpl = "pkg/sql/colexec/colexecwindow/window_peer_grouper_tmpl.go"

func genWindowPeerGrouperOps(inputFileContents string, wr io.Writer) error {
	s := strings.ReplaceAll(inputFileContents, "_PEER_GROUPER_STRING", "{{.String}}")

	// Now, generate the op, from the template.
	tmpl, err := template.New("peer_grouper_op").Parse(s)
	if err != nil {
		return err
	}

	windowPeerGrouperTmplInfos := []windowPeerGrouperTmplInfo{
		{AllPeers: false, HasPartition: false, String: "windowPeerGrouperNoPartition"},
		{AllPeers: false, HasPartition: true, String: "windowPeerGrouperWithPartition"},
		{AllPeers: true, HasPartition: false, String: "windowPeerGrouperAllPeersNoPartition"},
		{AllPeers: true, HasPartition: true, String: "windowPeerGrouperAllPeersWithPartition"},
	}
	return tmpl.Execute(wr, windowPeerGrouperTmplInfos)
}

func init() {
	registerGenerator(genWindowPeerGrouperOps, "window_peer_grouper.eg.go", windowPeerGrouperOpsTmpl)
}
