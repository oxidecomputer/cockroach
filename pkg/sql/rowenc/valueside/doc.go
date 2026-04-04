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

// Package valueside contains low-level primitives used to encode/decode SQL
// values into/from KV Values (see roachpb.Value).
//
// Low-level here means that these primitives do not operate with table or index
// descriptors.
//
// There are two separate schemes for encoding values:
//
//  - version 1 (legacy): the original encoding, which supported at most one SQL
//    value (column) per roachpb.Value. It is still used for old table
//    descriptors that went through many upgrades, and for some system tables.
//    Primitives related to this version contain the name `Legacy`.
//
//  - version 2 (column families): the current encoding which supports multiple
//    SQL values (columns) per roachpb.Value.
//
// See also: docs/tech-notes/encoding.md.
package valueside
