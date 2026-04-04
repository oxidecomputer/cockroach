// Copyright 2018 The Cockroach Authors.
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

package flagstub

import (
	"time"

	"github.com/cockroachdb/cockroach/pkg/roachprod/logger"
	"github.com/cockroachdb/cockroach/pkg/roachprod/vm"
	"github.com/cockroachdb/errors"
)

// New wraps a delegate vm.Provider to only return its name and
// flags.  This allows roachprod to provide a consistent tooling
// experience. Operations that can be reasonably stubbed out are
// implemented as no-op or no-value. All other operations will
// return the provided error.
func New(delegate vm.Provider, unimplemented string) vm.Provider {
	return &provider{delegate: delegate, unimplemented: unimplemented}
}

type provider struct {
	delegate      vm.Provider
	unimplemented string
}

// CleanSSH implements vm.Provider and is a no-op.
func (p *provider) CleanSSH() error {
	return nil
}

// ConfigSSH implements vm.Provider and is a no-op.
func (p *provider) ConfigSSH(zones []string) error {
	return nil
}

// Create implements vm.Provider and returns Unimplemented.
func (p *provider) Create(
	l *logger.Logger, names []string, opts vm.CreateOpts, providerOpts vm.ProviderOpts,
) error {
	return errors.Newf("%s", p.unimplemented)
}

// Delete implements vm.Provider and returns Unimplemented.
func (p *provider) Delete(vms vm.List) error {
	return errors.Newf("%s", p.unimplemented)
}

// Reset implements vm.Provider and is a no-op.
func (p *provider) Reset(vms vm.List) error {
	return nil
}

// Extend implements vm.Provider and returns Unimplemented.
func (p *provider) Extend(vms vm.List, lifetime time.Duration) error {
	return errors.Newf("%s", p.unimplemented)
}

// FindActiveAccount implements vm.Provider and returns an empty account.
func (p *provider) FindActiveAccount() (string, error) {
	return "", nil
}

// List implements vm.Provider and returns an empty list.
func (p *provider) List(l *logger.Logger) (vm.List, error) {
	return nil, nil
}

// Name implements vm.Provider and returns the delegate's name.
func (p *provider) Name() string {
	return p.delegate.Name()
}

// Active is part of the vm.Provider interface.
func (p *provider) Active() bool {
	return false
}

// ProjectActive is part of the vm.Provider interface.
func (p *provider) ProjectActive(project string) bool {
	return false
}

// CreateProviderFlags is part of the vm.Provider interface.
func (p *provider) CreateProviderOpts() vm.ProviderOpts {
	return nil
}
