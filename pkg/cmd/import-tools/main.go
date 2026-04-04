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

// import-tools adds a blank import to tools we use such that `go mod tidy`
// doesn't clean up needed dependencies when running `go install`.

//go:build tools
// +build tools

package main

import (
	"fmt"

	_ "github.com/aws/aws-sdk-go-v2"
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/cockroachdb/crlfmt"
	_ "github.com/cockroachdb/go-test-teamcity"
	_ "github.com/cockroachdb/gostdlib/cmd/gofmt"
	_ "github.com/cockroachdb/gostdlib/x/tools/cmd/goimports"
	_ "github.com/cockroachdb/stress"
	_ "github.com/cockroachdb/tools/cmd/stringer"
	_ "github.com/go-swagger/go-swagger/cmd/swagger"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/goware/modvendor"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway"
	_ "github.com/kevinburke/go-bindata/go-bindata"
	_ "github.com/kisielk/errcheck"
	_ "github.com/mattn/goveralls"
	_ "github.com/mibk/dupl"
	_ "github.com/mmatczuk/go_generics/cmd/go_generics"
	_ "github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc"
	_ "github.com/wadey/gocovmerge"
	_ "golang.org/x/perf/cmd/benchstat"
	_ "golang.org/x/tools/cmd/goyacc"
	_ "golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow"
	_ "honnef.co/go/tools/cmd/staticcheck"
)

func main() {
	fmt.Printf("You just lost the game\n")
}
