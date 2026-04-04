// Copyright 2015 The Cockroach Authors.
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

package config

import (
	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
)

// MakeZoneKeyPrefix returns the key prefix for id's row in the system.zones
// table.
func MakeZoneKeyPrefix(codec keys.SQLCodec, id descpb.ID) roachpb.Key {
	return codec.ZoneKeyPrefix(uint32(id))
}

// MakeZoneKey returns the key for a given id's entry in the system.zones table.
func MakeZoneKey(codec keys.SQLCodec, id descpb.ID) roachpb.Key {
	return codec.ZoneKey(uint32(id))
}

// DecodeObjectID decodes the object ID for the system-tenant from
// the front of key. It returns the decoded object ID, the remainder of the key,
// and whether the result is valid (i.e., whether the key was within the system
// tenant's structured key space).
func DecodeObjectID(codec keys.SQLCodec, key roachpb.RKey) (ObjectID, []byte, bool) {
	rem, id, err := codec.DecodeTablePrefix(key.AsRawKey())
	return ObjectID(id), rem, err == nil
}
