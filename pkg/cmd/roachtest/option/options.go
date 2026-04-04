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

package option

import (
	"github.com/cockroachdb/cockroach/pkg/roachprod"
	"github.com/cockroachdb/cockroach/pkg/roachprod/install"
)

// StartOpts is a type that combines the start options needed by roachprod and roachtest.
type StartOpts struct {
	RoachprodOpts install.StartOpts
	RoachtestOpts struct {
		Worker bool
	}
}

// DefaultStartOpts returns a StartOpts populated with default values.
func DefaultStartOpts() StartOpts {
	return StartOpts{RoachprodOpts: roachprod.DefaultStartOpts()}
}

// StopOpts is a type that combines the stop options needed by roachprod and roachtest.
type StopOpts struct {
	RoachprodOpts roachprod.StopOpts
	RoachtestOpts struct {
		Worker bool
	}
}

// DefaultStopOpts returns a StopOpts populated with default values.
func DefaultStopOpts() StopOpts {
	return StopOpts{RoachprodOpts: roachprod.DefaultStopOpts()}
}
