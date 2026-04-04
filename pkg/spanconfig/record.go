// Copyright 2022 The Cockroach Authors.
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

package spanconfig

import (
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/errors"
)

// Record ties a target to its corresponding config.
type Record struct {
	// target specifies the target (keyspan(s)) the config applies over.
	target Target

	// config is the set of attributes that apply over the corresponding target.
	config roachpb.SpanConfig
}

// MakeRecord returns a Record with the specified Target and SpanConfig. If the
// Record targets a SystemTarget, we also validate the SpanConfig.
func MakeRecord(target Target, cfg roachpb.SpanConfig) (Record, error) {
	if target.IsSystemTarget() {
		if err := cfg.ValidateSystemTargetSpanConfig(); err != nil {
			return Record{},
				errors.NewAssertionErrorWithWrappedErrf(err, "failed to validate SystemTarget SpanConfig")
		}
	}
	return Record{target: target, config: cfg}, nil
}

// IsEmpty returns true if the receiver is an empty Record.
func (r *Record) IsEmpty() bool {
	return r.target.isEmpty() && r.config.IsEmpty()
}

// GetTarget returns the Record target.
func (r *Record) GetTarget() Target {
	return r.target
}

// GetConfig returns the Record config.
func (r *Record) GetConfig() roachpb.SpanConfig {
	return r.config
}
