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

package settingswatcher

import "github.com/cockroachdb/cockroach/pkg/settings"

// OverridesMonitor is an interface through which the settings watcher can
// receive setting overrides. Used for non-system tenants.
//
// The expected usage is to listen for a message on NotifyCh(), and use
// Current() to retrieve the updated list of overrides when a message is
// received.
type OverridesMonitor interface {
	// RegisterOverridesChannel returns a channel that receives a message
	// any time the current set of overrides changes.
	// The channel receives an initial event immediately.
	RegisterOverridesChannel() <-chan struct{}

	// Overrides retrieves the current set of setting overrides, as a map from
	// setting key to EncodedValue. Any settings that are present must be set to
	// the overridden value.
	Overrides() map[string]settings.EncodedValue
}
