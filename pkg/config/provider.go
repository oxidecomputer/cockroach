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

package config

// SystemConfigProvider is capable of providing the SystemConfig, as well as
// notifying clients of updates to the SystemConfig.
type SystemConfigProvider interface {
	// GetSystemConfig returns the local unmarshaled version of the system
	// config. Returns nil if the system config hasn't been set yet.
	GetSystemConfig() *SystemConfig

	// RegisterSystemConfigChannel registers a channel to signify updates for
	// the system config. It is notified after registration (if a system config
	// is already set), and whenever a new system config is successfully
	// unmarshaled.
	RegisterSystemConfigChannel() (_ <-chan struct{}, unregister func())
}

// ConstantSystemConfigProvider is an implementation of SystemConfigProvider which
// always returns the same value.
type ConstantSystemConfigProvider struct {
	cfg *SystemConfig
}

// NewConstantSystemConfigProvider constructs a SystemConfigProvider which
// always returns the same value.
func NewConstantSystemConfigProvider(cfg *SystemConfig) *ConstantSystemConfigProvider {
	p := &ConstantSystemConfigProvider{cfg: cfg}
	return p
}

// GetSystemConfig implements the SystemConfigProvider interface.
func (c *ConstantSystemConfigProvider) GetSystemConfig() *SystemConfig {
	return c.cfg
}

// RegisterSystemConfigChannel implements the SystemConfigProvider interface.
func (c *ConstantSystemConfigProvider) RegisterSystemConfigChannel() (
	_ <-chan struct{},
	unregister func(),
) {
	// The system config will never be updated, so return a nil channel.
	return nil, func() {}
}
