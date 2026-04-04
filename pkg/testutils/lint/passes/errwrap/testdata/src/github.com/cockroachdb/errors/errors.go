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

package errors

import "fmt"

func New(msg string) error {
	return fmt.Errorf(msg)
}

func Newf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func Wrap(_ error, msg string) error {
	return fmt.Errorf(msg)
}

func Wrapf(_ error, format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func WrapWithDepth(depth int, err error, msg string) error {
	return fmt.Errorf(msg)
}

func WrapWithDepthf(depth int, err error, format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func NewWithDepth(_ int, msg string) error {
	return fmt.Errorf(msg)
}

func NewWithDepthf(_ int, format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func AssertionFailedf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func AssertionFailedWithDepthf(_ int, format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func NewAssertionErrorWithWrappedErrf(_ error, format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
