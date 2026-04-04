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

package clierrorplus

import (
	"github.com/cockroachdb/cockroach/pkg/cli/clierror"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/spf13/cobra"
)

// MaybeShoutError calls log.Shout on errors, better surfacing problems to the user.
func MaybeShoutError(
	wrapped func(*cobra.Command, []string) error,
) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		err := wrapped(cmd, args)
		return CheckAndMaybeShout(err)
	}
}

// CheckAndMaybeShout shouts the error, if non-nil to the OPS logging
// channel.
func CheckAndMaybeShout(err error) error {
	return clierror.CheckAndMaybeLog(err, log.Ops.Shoutf)
}
