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

package a

import "github.com/cockroachdb/cockroach/pkg/util"

type onlyField struct {
	_ util.NoCopy
}

type firstField struct {
	_ util.NoCopy
	a int64
}

type middleField struct {
	a int64
	_ util.NoCopy // want `Illegal use of util.NoCopy - must be first field in struct`
	b int64
}

type lastField struct {
	a int64
	_ util.NoCopy // want `Illegal use of util.NoCopy - must be first field in struct`
}

type embeddedField struct {
	util.NoCopy // want `Illegal use of util.NoCopy - should not be embedded`
}

type multiField struct {
	_, _ util.NoCopy // want `Illegal use of util.NoCopy - should be included only once`
}

type namedField struct {
	noCopy util.NoCopy // want `Illegal use of util.NoCopy - should be unnamed`
}
