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

package securitytest

import (
	"os"
	"path/filepath"

	"github.com/cockroachdb/cockroach/pkg/security"
)

// CreateTestCerts populates the test certificates in the given directory.
func CreateTestCerts(certsDir string) (cleanup func() error) {
	// Copy these assets to disk from embedded strings, so this test can
	// run from a standalone binary.
	// Disable embedded certs, or the security library will try to load
	// our real files as embedded assets.
	security.ResetAssetLoader()

	assets := []string{
		filepath.Join(security.EmbeddedCertsDir, security.EmbeddedCACert),
		filepath.Join(security.EmbeddedCertsDir, security.EmbeddedCAKey),
		filepath.Join(security.EmbeddedCertsDir, security.EmbeddedNodeCert),
		filepath.Join(security.EmbeddedCertsDir, security.EmbeddedNodeKey),
		filepath.Join(security.EmbeddedCertsDir, security.EmbeddedRootCert),
		filepath.Join(security.EmbeddedCertsDir, security.EmbeddedRootKey),
		filepath.Join(security.EmbeddedCertsDir, security.EmbeddedTenantCACert),
	}

	for _, a := range assets {
		_, err := RestrictedCopy(a, certsDir, filepath.Base(a))
		if err != nil {
			panic(err)
		}
	}

	return func() error {
		security.SetAssetLoader(EmbeddedAssets)
		return os.RemoveAll(certsDir)
	}
}
