// Copyright 2019 The Cockroach Authors.
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

package reports

import (
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/config"
)

// ZoneKey is the index of the first level in the constraint conformance report.
type ZoneKey struct {
	// ZoneID is the id of the zone this key is referencing.
	ZoneID config.ObjectID
	// SubzoneID identifies what subzone, if any, this key is referencing. The
	// zero value (also named NoSubzone) indicates that the key is referring to a
	// zone, not a subzone.
	SubzoneID base.SubzoneID
}

// NoSubzone is used inside a zoneKey to indicate that the key represents a
// zone, not a subzone.
const NoSubzone base.SubzoneID = 0

// MakeZoneKey creates a zoneKey.
//
// Use NoSubzone for subzoneID to indicate that the key references a zone, not a
// subzone.
func MakeZoneKey(zoneID config.ObjectID, subzoneID base.SubzoneID) ZoneKey {
	return ZoneKey{
		ZoneID:    zoneID,
		SubzoneID: subzoneID,
	}
}

func (k ZoneKey) String() string {
	return fmt.Sprintf("%d,%d", k.ZoneID, k.SubzoneID)
}

// Less compares two ZoneKeys.
func (k ZoneKey) Less(other ZoneKey) bool {
	if k.ZoneID < other.ZoneID {
		return true
	}
	if k.ZoneID > other.ZoneID {
		return false
	}
	return k.SubzoneID < other.SubzoneID
}
