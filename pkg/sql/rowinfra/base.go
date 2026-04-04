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

// Package rowinfra contains constants and types used by the row package
// that must also be accessible from other packages.
package rowinfra

// RowLimit represents a response limit expressed in terms of number of result
// rows. RowLimits get ultimately converted to KeyLimits and are translated into
// BatchRequest.MaxSpanRequestKeys.
type RowLimit uint64

// KeyLimit represents a response limit expressed in terms of number of keys.
type KeyLimit int64

// BytesLimit represents a response limit expressed in terms of the size of the
// results. A BytesLimit ultimately translates into BatchRequest.TargetBytes.
type BytesLimit uint64

// NoRowLimit can be passed to Fetcher.StartScan to signify that the caller
// doesn't want to limit the number of result rows for each scan request.
const NoRowLimit RowLimit = 0

// NoBytesLimit can be passed to Fetcher.StartScan to signify that the caller
// doesn't want to limit the size of results for each scan request.
//
// See also DefaultBatchBytesLimit.
const NoBytesLimit BytesLimit = 0

// ProductionKVBatchSize is the kv batch size to use for production (i.e.,
// non-test) clusters.
const ProductionKVBatchSize KeyLimit = 100000

// DefaultBatchBytesLimit is the maximum number of bytes a scan request can
// return.
const DefaultBatchBytesLimit BytesLimit = 10 << 20 // 10 MB
