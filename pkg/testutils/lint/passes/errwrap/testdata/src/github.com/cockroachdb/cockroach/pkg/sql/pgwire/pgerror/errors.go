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

package pgerror

import (
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
)

func NewWithDepthf(depth int, code pgcode.Code, format string, args ...interface{}) error {
	return fmt.Errorf(format, args)
}

func New(code pgcode.Code, msg string) error {
	return fmt.Errorf(msg)

}

func Newf(code pgcode.Code, format string, args ...interface{}) error {
	return fmt.Errorf(format, args)
}

func Wrapf(err error, code pgcode.Code, format string, args ...interface{}) error {
	return fmt.Errorf(format, args)
}

func WrapWithDepthf(
	depth int, err error, code pgcode.Code, format string, args ...interface{},
) error {
	return fmt.Errorf(format, args)

}

func Wrap(err error, code pgcode.Code, msg string) error {
	return fmt.Errorf(msg)
}
