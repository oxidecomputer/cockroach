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

package democluster

import (
	"strings"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
)

// DemoLocalityList represents a list of localities for the cockroach
// demo command.
type DemoLocalityList []roachpb.Locality

// Type implements the pflag.Value interface.
func (l *DemoLocalityList) Type() string { return "demoLocalityList" }

// String implements the pflag.Value interface.
func (l *DemoLocalityList) String() string {
	s := ""
	for _, loc := range []roachpb.Locality(*l) {
		s += loc.String()
	}
	return s
}

// Set implements the pflag.Value interface.
func (l *DemoLocalityList) Set(value string) error {
	*l = []roachpb.Locality{}
	locs := strings.Split(value, ":")
	for _, value := range locs {
		parsedLoc := &roachpb.Locality{}
		if err := parsedLoc.Set(value); err != nil {
			return err
		}
		*l = append(*l, *parsedLoc)
	}
	return nil
}

var defaultLocalities = DemoLocalityList{
	// Default localities for a 3 node cluster
	{Tiers: []roachpb.Tier{{Key: "region", Value: "us-east1"}, {Key: "az", Value: "b"}}},
	{Tiers: []roachpb.Tier{{Key: "region", Value: "us-east1"}, {Key: "az", Value: "c"}}},
	{Tiers: []roachpb.Tier{{Key: "region", Value: "us-east1"}, {Key: "az", Value: "d"}}},
	// Default localities for a 6 node cluster
	{Tiers: []roachpb.Tier{{Key: "region", Value: "us-west1"}, {Key: "az", Value: "a"}}},
	{Tiers: []roachpb.Tier{{Key: "region", Value: "us-west1"}, {Key: "az", Value: "b"}}},
	{Tiers: []roachpb.Tier{{Key: "region", Value: "us-west1"}, {Key: "az", Value: "c"}}},
	// Default localities for a 9 node cluster
	{Tiers: []roachpb.Tier{{Key: "region", Value: "europe-west1"}, {Key: "az", Value: "b"}}},
	{Tiers: []roachpb.Tier{{Key: "region", Value: "europe-west1"}, {Key: "az", Value: "c"}}},
	{Tiers: []roachpb.Tier{{Key: "region", Value: "europe-west1"}, {Key: "az", Value: "d"}}},
}
