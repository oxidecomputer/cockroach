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

package catpb

import (
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/errors"
)

// RegionName is an alias for a region stored on the database.
type RegionName string

// String implements fmt.Stringer.
func (r RegionName) String() string {
	return string(r)
}

// RegionNames is an alias for a slice of regions.
type RegionNames []RegionName

// ToStrings converts the RegionNames slice to a string slice.
func (regions RegionNames) ToStrings() []string {
	ret := make([]string, len(regions))
	for i, region := range regions {
		ret[i] = string(region)
	}
	return ret
}

// TelemetryName returns the name to use for the given locality.
func (cfg *LocalityConfig) TelemetryName() (string, error) {
	switch l := cfg.Locality.(type) {
	case *LocalityConfig_Global_:
		return tree.TelemetryNameGlobal, nil
	case *LocalityConfig_RegionalByTable_:
		if l.RegionalByTable.Region != nil {
			return tree.TelemetryNameRegionalByTableIn, nil
		}
		return tree.TelemetryNameRegionalByTable, nil
	case *LocalityConfig_RegionalByRow_:
		if l.RegionalByRow.As != nil {
			return tree.TelemetryNameRegionalByRowAs, nil
		}
		return tree.TelemetryNameRegionalByRow, nil
	}
	return "", errors.AssertionFailedf(
		"unknown locality config TelemetryName: type %T",
		cfg.Locality,
	)
}
