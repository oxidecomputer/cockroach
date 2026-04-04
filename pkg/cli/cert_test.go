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

package cli

func Example_cert() {
	c := NewCLITest(TestCLIParams{})
	defer c.Cleanup()

	c.RunWithCAArgs([]string{"cert", "create-client", "foo"})
	c.RunWithCAArgs([]string{"cert", "create-client", "Ομηρος"})
	c.RunWithCAArgs([]string{"cert", "create-client", "0foo"})
	c.RunWithCAArgs([]string{"cert", "create-client", ",foo"})

	// Output:
	// cert create-client foo
	// cert create-client Ομηρος
	// cert create-client 0foo
	// cert create-client ,foo
	// ERROR: failed to generate client certificate and key: username is invalid
	// HINT: Usernames are case insensitive, must start with a letter, digit or underscore, may contain letters, digits, dashes, periods, or underscores, and must not exceed 63 characters.
}
