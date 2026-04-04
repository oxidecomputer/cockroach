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

package rangefeed

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
)

// NewDBAdapter allows tests to construct a dbAdapter.
var NewDBAdapter = newDBAdapter

// NewFactoryWithDB allows tests to construct a factory with an injected db.
var NewFactoryWithDB = newFactory

// KVDB forwards the definition of DB to tests.
type KVDB = DB

// ScanConfig forwards the definition of scanConfig to tests.
type ScanConfig = scanConfig

// ScanWithOptions is exposed for testing in order to call Scan with scanConfig
// extracted from the specified list of options.
func (dbc *dbAdapter) ScanWithOptions(
	ctx context.Context,
	spans []roachpb.Span,
	asOf hlc.Timestamp,
	rowFn func(value roachpb.KeyValue),
	opts ...Option,
) error {
	var c config
	initConfig(&c, opts)
	return dbc.Scan(ctx, spans, asOf, rowFn, c.scanConfig)
}
