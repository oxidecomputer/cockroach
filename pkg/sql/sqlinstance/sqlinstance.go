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

// Package sqlinstance provides interfaces that will be exposed
// to interact with the sqlinstance subsystem.
// This subsystem will initialize and maintain a unique instance ID
// per SQL instance along with mapping of an active instance ID to
// its SQL address.
package sqlinstance

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlliveness"
	"github.com/cockroachdb/errors"
)

// InstanceInfo exposes information on a SQL instance such as ID, network address and
// the associated sqlliveness.SessionID.
type InstanceInfo struct {
	InstanceID   base.SQLInstanceID
	InstanceAddr string
	SessionID    sqlliveness.SessionID
}

// AddressResolver exposes API for retrieving the instance address and all live instances for a tenant.
type AddressResolver interface {
	// GetInstance returns the InstanceInfo for a pod given the instance ID.
	// Returns a NonExistentInstanceError if an instance with the given ID
	// does not exist.
	GetInstance(context.Context, base.SQLInstanceID) (InstanceInfo, error)
	// GetAllInstances returns a list of all SQL instances for the tenant.
	GetAllInstances(context.Context) ([]InstanceInfo, error)
}

// Provider is a wrapper around sqlinstance subsystem for external consumption.
type Provider interface {
	AddressResolver
	// Instance returns the instance ID and sqlliveness.SessionID for the
	// current SQL instance.
	Instance(context.Context) (base.SQLInstanceID, sqlliveness.SessionID, error)
	// Start starts the instanceprovider. This will block until
	// the underlying instance data reader has been started.
	Start(context.Context) error
}

// NonExistentInstanceError can be returned if a SQL instance does not exist.
var NonExistentInstanceError = errors.Errorf("non existent SQL instance")

// NotStartedError can be returned if the sqlinstance subsystem has not been started yet.
var NotStartedError = errors.Errorf("sqlinstance subsystem not started")
