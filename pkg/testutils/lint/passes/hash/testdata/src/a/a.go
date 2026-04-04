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

package a

import "crypto/sha256"

func init() {
	var inputBytes, hashedBytes []byte
	_ = hashedBytes

	{
		h := sha256.New()
		h.Write(inputBytes)
		hashedBytes = h.Sum(nil)
	}

	{
		h := sha256.New()
		h.Write(inputBytes)
		var hashedBytes [sha256.Size]byte
		h.Sum(hashedBytes[:0])
	}

	{
		hashedBytes = sha256.New().Sum(inputBytes) // want `probable misuse of hash.Hash.Sum: provide parameter or use return value, but not both`
	}
}
