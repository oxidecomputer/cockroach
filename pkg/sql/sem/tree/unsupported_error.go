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

package tree

var _ error = &UnsupportedError{}

// UnsupportedError is an error object which is returned by some unimplemented SQL
// statements. It is currently only used to skip over PGDUMP statements during
// an import.
type UnsupportedError struct {
	Err         error
	FeatureName string
}

func (u *UnsupportedError) Error() string {
	return u.Err.Error()
}

// Cause implements causer.
func (u *UnsupportedError) Cause() error { return u.Err }

// Unwrap implements wrapper.
func (u *UnsupportedError) Unwrap() error { return u.Err }
