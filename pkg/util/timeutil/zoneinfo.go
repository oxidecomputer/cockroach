// Copyright 2016 The Cockroach Authors.
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

package timeutil

import (
	"strings"
	"time"
	// embed tzdata in case system tzdata is not available.
	_ "time/tzdata"
)

//go:generate go run gen/main.go

// LoadLocation returns the time.Location with the given name.
// The name is taken to be a location name corresponding to a file
// in the IANA Time Zone database, such as "America/New_York".
//
// We do not use Go's time.LoadLocation() directly because it maps
// "Local" to the local time zone, whereas we want UTC.
func LoadLocation(name string) (*time.Location, error) {
	loweredName := strings.ToLower(name)
	switch loweredName {
	case "local", "default":
		loweredName = "utc"
		name = "UTC"
	}
	// If we know this is a lowercase name in tzdata, use the uppercase form.
	if v, ok := lowercaseTimezones[loweredName]; ok {
		// If this location is not found, we may have a case where the tzdata names
		// have different values than the system tz names.
		// If this is the case, allback onto the default logic, where the name is read
		// off other sources before tzdata.
		if loc, err := time.LoadLocation(v); err == nil {
			return loc, nil
		}
	}
	return time.LoadLocation(name)
}
