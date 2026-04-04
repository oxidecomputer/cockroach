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

package rel

import (
	"reflect"

	"github.com/cockroachdb/errors"
)

var (
	reflectTypeType    = reflect.TypeOf((*reflect.Type)(nil)).Elem()
	emptyInterfaceType = reflect.TypeOf((*interface{})(nil)).Elem()
)

func makeComparableValue(val interface{}) (typedValue, error) {
	if typ, isType := val.(reflect.Type); isType {
		return typedValue{
			typ:   reflectTypeType,
			value: typ,
		}, nil
	}
	vv := reflect.ValueOf(val)
	if err := checkNotNil(vv); err != nil {
		return typedValue{}, err
	}
	typ := vv.Type()
	switch {
	case isSupportScalarKind(typ.Kind()):
		// We need to allocate a new pointer.
		compType := getComparableType(typ)
		vvNew := reflect.New(vv.Type())
		vvNew.Elem().Set(vv)
		return typedValue{
			typ:   vv.Type(),
			value: vvNew.Convert(reflect.PtrTo(compType)).Interface(),
		}, nil
	case typ.Kind() == reflect.Ptr:
		switch {
		case isSupportScalarKind(typ.Elem().Kind()):
			compType := getComparableType(typ.Elem())
			return typedValue{
				typ:   vv.Type().Elem(),
				value: vv.Convert(reflect.PtrTo(compType)).Interface(),
			}, nil
		case typ.Elem().Kind() == reflect.Struct:
			return typedValue{
				typ:   vv.Type(),
				value: val,
			}, nil
		default:
			return typedValue{}, errors.Errorf(
				"unsupported pointer kind %v for type %T", typ.Elem().Kind(), val,
			)
		}
	default:
		return typedValue{}, errors.Errorf(
			"unsupported kind %v for type %T", typ.Kind(), val,
		)
	}
}

func checkNotNil(v reflect.Value) error {
	if !v.IsValid() {
		// you are not allowed to put A nil pointer here
		return errors.Errorf("invalid nil")
	}
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return errors.Errorf("invalid nil %v", v.Type())
	}
	return nil
}
