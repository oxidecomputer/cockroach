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

//go:build generator
// +build generator

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/cockroachdb/cockroach/pkg/cli/exit"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scpb"
)

var (
	out = flag.String("out", "", "output file for generated UML")
)

func main() {
	flag.Parse()
	if err := run(*out); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		exit.WithCode(exit.FatalError())
	}
}
func run(out string) error {
	if out == "" {
		return fmt.Errorf("output required")
	}
	var buf bytes.Buffer
	var parentRelations bytes.Buffer

	getParentsFromField := func(f reflect.StructField) []string {
		if parentTag := f.Tag.Get("parent"); parentTag != "" {
			return strings.Split(parentTag, ", ")
		}
		return nil
	}

	buf.WriteString("@startuml\n")
	elementProtoType := reflect.TypeOf((*scpb.ElementProto)(nil)).Elem()
	for i := 0; i < elementProtoType.NumField(); i++ {
		fieldType := elementProtoType.Field(i).Type.Elem()
		buf.WriteString(fmt.Sprintf(
			"object %s\n\n",
			fieldType.Name()))
		for j := 0; j < fieldType.NumField(); j++ {
			arrayPrefix := " "
			if fieldType.Field(j).Type.Kind() == reflect.Slice {
				arrayPrefix = "[]"
			}
			buf.WriteString(
				fmt.Sprintf("%s : %s%s\n",
					fieldType.Name(),
					arrayPrefix,
					fieldType.Field(j).Name),
			)
		}
		buf.WriteString("\n")
		// The parent tag has a list of elements that are the parents
		// to this element. We will collect these and emit them later
		// in the PlantUML syntax.
		for _, parent := range getParentsFromField(elementProtoType.Field(i)) {
			parentRelations.WriteString(fmt.Sprintf(
				"%s <|-- %s\n", parent, fieldType.Name()))
		}
	}
	// Append all the object relationships at
	// the end.
	buf.Write(parentRelations.Bytes())
	buf.WriteString("@enduml\n")
	return ioutil.WriteFile(out, buf.Bytes(), 0777)
}
