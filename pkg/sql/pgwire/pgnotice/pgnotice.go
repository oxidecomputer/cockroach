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

package pgnotice

import (
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/errors"
)

// Notice is an wrapper around errors that are intended to be notices.
type Notice error

// Newf generates a Notice with a format string.
func Newf(format string, args ...interface{}) Notice {
	err := errors.NewWithDepthf(1, format, args...)
	err = pgerror.WithCandidateCode(err, pgcode.SuccessfulCompletion)
	err = pgerror.WithSeverity(err, "NOTICE")
	return Notice(err)
}

// NewWithSeverityf generates a Notice with a format string and severity.
func NewWithSeverityf(severity string, format string, args ...interface{}) Notice {
	err := errors.NewWithDepthf(1, format, args...)
	err = pgerror.WithCandidateCode(err, pgcode.SuccessfulCompletion)
	err = pgerror.WithSeverity(err, severity)
	return Notice(err)
}
