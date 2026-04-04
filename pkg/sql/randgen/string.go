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

package randgen

import "math/rand"

// RandString generates a random string of the desired length from the
// input alphabet.
func RandString(rng *rand.Rand, length int, alphabet string) string {
	buf := make([]byte, length)
	for i := range buf {
		buf[i] = alphabet[rng.Intn(len(alphabet))]
	}
	return string(buf)
}
