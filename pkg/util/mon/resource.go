// Copyright 2017 The Cockroach Authors.
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

package mon

import (
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/errors"
)

// Resource is an interface used to abstract the specifics of tracking bytes
// usage by different types of resources.
type Resource interface {
	NewBudgetExceededError(requestedBytes int64, reservedBytes int64, budgetBytes int64) error
}

func newMemoryBudgetExceededError(
	requestedBytes int64, reservedBytes int64, budgetBytes int64,
) error {
	return pgerror.WithCandidateCode(
		errors.Newf(
			"memory budget exceeded: %d bytes requested, %d currently allocated, %d bytes in budget",
			errors.Safe(requestedBytes),
			errors.Safe(reservedBytes),
			errors.Safe(budgetBytes),
		),
		pgcode.OutOfMemory,
	)
}

// memoryResource is a Resource that represents memory.
type memoryResource struct{}

// MemoryResource is a utility singleton used as an argument when creating a
// BytesMonitor to indicate that the monitor will be tracking memory usage.
var MemoryResource Resource = memoryResource{}

// NewBudgetExceededError implements the Resource interface.
func (m memoryResource) NewBudgetExceededError(
	requestedBytes int64, reservedBytes int64, budgetBytes int64,
) error {
	return newMemoryBudgetExceededError(requestedBytes, reservedBytes, budgetBytes)
}

// memoryResourceWithErrorHint is a Resource that represents memory and augments
// the "budget exceeded" error with a hint.
type memoryResourceWithErrorHint struct {
	hint string
}

// NewMemoryResourceWithErrorHint returns a new memory Resource that augments
// all "budget exceeded" errors with the given hint.
func NewMemoryResourceWithErrorHint(hint string) Resource {
	return memoryResourceWithErrorHint{hint: hint}
}

// NewBudgetExceededError implements the Resource interface.
func (m memoryResourceWithErrorHint) NewBudgetExceededError(
	requestedBytes int64, reservedBytes int64, budgetBytes int64,
) error {
	return errors.WithHint(
		newMemoryBudgetExceededError(requestedBytes, reservedBytes, budgetBytes),
		m.hint,
	)
}

// diskResource is a Resource that represents disk.
type diskResource struct{}

// DiskResource is a utility singleton used as an argument when creating a
// BytesMonitor to indicate that the monitor will be tracking disk usage.
var DiskResource Resource = diskResource{}

// NewBudgetExceededError implements the Resource interface.
func (d diskResource) NewBudgetExceededError(
	requestedBytes int64, reservedBytes int64, budgetBytes int64,
) error {
	return pgerror.WithCandidateCode(
		errors.Newf(
			"disk budget exceeded: %d bytes requested, %d currently allocated, %d bytes in budget",
			errors.Safe(requestedBytes),
			errors.Safe(reservedBytes),
			errors.Safe(budgetBytes),
		), pgcode.DiskFull)
}
