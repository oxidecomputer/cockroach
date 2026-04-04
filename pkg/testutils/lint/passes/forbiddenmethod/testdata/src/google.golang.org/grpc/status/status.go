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

package status

import (
	istatus "google.golang.org/grpc/internal/status"
)

// This mirrors the gRPC structure, where status.Status is an alias for an
// internal type. The lint needs to check on the internal type.
type Status = istatus.Status

func New(_ int, _ string) *Status {
	return &Status{}
}
