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

package tests

import (
	// Import embed for the version map
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/util/version"
	"github.com/cockroachdb/errors"
)

// You can update the values in predecessor_version.json to point at newer patch releases.
//
// NB: If a new key was added (e.g. if you're releasing a new major
// release), you'll also need to regenerate fixtures. To regenerate
// fixtures, you will need to run acceptance/version-upgrade with the
// checkpoint option enabled to create the missing store directory
// fixture (see runVersionUpgrade).
//go:embed predecessor_version.json
var verMapJSON []byte

type versionMap map[string]string

func unmarshalVersionMap() (versionMap, error) {
	var res versionMap
	err := json.Unmarshal(verMapJSON, &res)
	if err != nil {
		return versionMap{}, err
	}
	return res, nil
}

// PredecessorVersion returns a recent predecessor of the build version (i.e.
// the build tag of the main binary). For example, if the running binary is from
// the master branch prior to releasing 19.2.0, this will return a recent
// (ideally though not necessarily the latest) 19.1 patch release.
func PredecessorVersion(buildVersion version.Version) (string, error) {
	if buildVersion == (version.Version{}) {
		return "", errors.Errorf("buildVersion not set")
	}

	verMap, err := unmarshalVersionMap()
	if err != nil {
		return "", errors.Wrap(err, "cannot load version map")
	}
	buildVersionMajorMinor := fmt.Sprintf("%d.%d", buildVersion.Major(), buildVersion.Minor())
	v, ok := verMap[buildVersionMajorMinor]
	if !ok {
		return "", errors.Errorf("prev version not set for version: %s", buildVersionMajorMinor)
	}
	return v, nil
}
