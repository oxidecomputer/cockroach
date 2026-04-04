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

package opgen

import (
	"reflect"
	"sort"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scpb"
)

func TestOpGen(t *testing.T) {
	var elementProto scpb.ElementProto
	elementProtoType := reflect.ValueOf(elementProto).Type()
	fieldIndexes := make([]int, elementProtoType.NumField())
	for i, n := 0, elementProtoType.NumField(); i < n; i++ {
		fieldIndexes[i] = i
	}
	sort.Slice(fieldIndexes, func(i, j int) bool {
		return elementProtoType.Field(fieldIndexes[i]).Name < elementProtoType.Field(fieldIndexes[j]).Name
	})
	for _, i := range fieldIndexes {
		field := elementProtoType.Field(i)
		t.Run(field.Name, func(t *testing.T) {
			var adds, drops []target
			for _, tg := range opRegistry.targets {
				if reflect.ValueOf(tg.e).Type() == field.Type {
					switch tg.status {
					case scpb.Status_PUBLIC:
						adds = append(adds, tg)
					case scpb.Status_ABSENT:
						drops = append(drops, tg)
					}
				}
			}
			if len(adds) != 1 {
				t.Errorf("expected one registered adding spec for %s, instead found %d", field.Name, len(adds))
			}
			if len(drops) != 1 {
				t.Errorf("expected one registered dropping spec for %s, instead found %d", field.Name, len(drops))
			}
		})
	}
}
