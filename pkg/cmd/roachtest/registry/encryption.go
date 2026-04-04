// Copyright 2022 The Cockroach Authors.
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

package registry

import "fmt"

// EncryptionSupport encodes the relationship of a test with
// encryption-at-rest. Tests can either opt-in to metamorphic
// encryption, or require that encryption is always on or always off
// (default).
type EncryptionSupport int

func (es EncryptionSupport) String() string {
	switch es {
	case EncryptionAlwaysEnabled:
		return "always-enabled"
	case EncryptionAlwaysDisabled:
		return "always-disabled"
	case EncryptionMetamorphic:
		return "metamorphic"
	default:
		return fmt.Sprintf("unknown-%d", es)
	}
}

const (
	// EncryptionAlwaysDisabled indicates that the test requires
	// encryption to be disabled. The test will only run on clusters
	// with encryption disabled.
	EncryptionAlwaysDisabled = EncryptionSupport(iota)
	// EncryptionAlwaysEnabled indicates that the test requires
	// encryption to be enabled. The test will only run on clusters
	// with encryption enabled.
	EncryptionAlwaysEnabled
	// EncryptionMetamorphic indicates that a test opted-in to
	// metamorphic encryption. Whether the test runs on a cluster with
	// encryption enabled depends on the probability passed to
	// --metamorphic-encryption-probability.
	EncryptionMetamorphic
)
