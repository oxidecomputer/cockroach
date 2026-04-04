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

import "sync"

// valuesMap is a container for attributes.
//
// It stores the data in a format which is convenient for performing
// comparisons and lookups. If you want strongly typed data out of it,
// you need to use a Schema to retrieve that data. Note that the library
// expects all values to be stored in the map in the comparable, primitive
// form and not in the strongly typed format.
type valuesMap struct {
	attrs ordinalSet
	m     map[ordinal]interface{}
}

var valuesSyncPool = sync.Pool{
	New: func() interface{} {
		return &valuesMap{
			m: make(map[ordinal]interface{}),
		}
	},
}

func getValues() *valuesMap {
	return valuesSyncPool.Get().(*valuesMap)
}

func putValues(v *valuesMap) {
	v.clear()
	valuesSyncPool.Put(v)
}

func (vm *valuesMap) clear() {
	for k := range vm.m {
		delete(vm.m, k)
	}
	vm.attrs = 0
}

// get retrieves the primitive valuesMap stores in the valuesMap
// struct.
func (vm valuesMap) get(a ordinal) interface{} {
	return vm.m[a]
}

func (vm *valuesMap) add(ord ordinal, v interface{}) {
	vm.attrs = vm.attrs.add(ord)
	vm.m[ord] = v
}
