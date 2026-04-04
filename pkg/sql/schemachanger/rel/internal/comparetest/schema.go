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

package comparetest

import (
	"reflect"

	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/rel"
)

type entity struct {
	I8       int8
	PI8      *int8
	I16      int16
	PI16     *int16
	I32      int32
	PI32     *int32
	I64      int64
	PI64     *int64
	UI8      uint8
	PUI8     *uint8
	UI16     uint16
	PUI16    *uint16
	UI32     uint32
	PUI32    *uint32
	UI64     uint64
	PUI64    *uint64
	I        int
	PI       *int
	UI       uint
	PUI      *uint
	Str      string
	PStr     *string
	Uintptr  uintptr
	PUintptr *uintptr
}

// testAttr is a rel.Attr used for testing.
type testAttr int8

var _ rel.Attr = testAttr(0)

//go:generate stringer --type testAttr  --tags test
const (
	i8 testAttr = iota
	pi8
	i16
	pi16
	i32
	pi32
	i64
	pi64
	ui8
	pui8
	ui16
	pui16
	ui32
	pui32
	ui64
	pui64
	i
	pi
	ui
	pui
	str
	pstr
	_uintptr
	puintptr
)

var schema = rel.MustSchema(
	"testschema",
	rel.EntityMapping(reflect.TypeOf((*entity)(nil)),
		rel.EntityAttr(i8, "I8"),
		rel.EntityAttr(pi8, "PI8"),
		rel.EntityAttr(i16, "I16"),
		rel.EntityAttr(pi16, "PI16"),
		rel.EntityAttr(i32, "I32"),
		rel.EntityAttr(pi32, "PI32"),
		rel.EntityAttr(i64, "I64"),
		rel.EntityAttr(pi64, "PI64"),
		rel.EntityAttr(ui8, "UI8"),
		rel.EntityAttr(pui8, "PUI8"),
		rel.EntityAttr(ui16, "UI16"),
		rel.EntityAttr(pui16, "PUI16"),
		rel.EntityAttr(ui32, "UI32"),
		rel.EntityAttr(pui32, "PUI32"),
		rel.EntityAttr(ui64, "UI64"),
		rel.EntityAttr(pui64, "PUI64"),
		rel.EntityAttr(i, "I"),
		rel.EntityAttr(pi, "PI"),
		rel.EntityAttr(ui, "UI"),
		rel.EntityAttr(pui, "PUI"),
		rel.EntityAttr(str, "Str"),
		rel.EntityAttr(pstr, "PStr"),
		rel.EntityAttr(_uintptr, "Uintptr"),
		rel.EntityAttr(puintptr, "PUintptr"),
	))
