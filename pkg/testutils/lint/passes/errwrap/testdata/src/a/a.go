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

package a

import (
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/errors"
)

var wrappedErr = fmt.Errorf("abc")
var anotherErr = fmt.Errorf("def")

func init() {
	_ = recover()

	_ = fmt.Errorf(wrappedErr.Error())              // want `err.Error\(\) is passed to fmt.Errorf; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`
	_ = fmt.Errorf("format %s", wrappedErr.Error()) // want `err.Error\(\) is passed to fmt.Errorf; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`

	s := wrappedErr.Error()
	_ = fmt.Errorf("format %s", s) // this way is allowed

	_ = pgerror.Wrap(anotherErr, pgcode.Warning, wrappedErr.Error())                           // want `err.Error\(\) is passed to pgerror.Wrap; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`
	_ = pgerror.Wrapf(anotherErr, pgcode.Warning, "format %s", wrappedErr.Error())             // want `err.Error\(\) is passed to pgerror.Wrapf; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`
	_ = pgerror.WrapWithDepthf(1, anotherErr, pgcode.Warning, "format %s", wrappedErr.Error()) // want `err.Error\(\) is passed to pgerror.WrapWithDepthf; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`
	_ = pgerror.New(pgcode.Warning, wrappedErr.Error())                                        // want `err.Error\(\) is passed to pgerror.New; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`
	_ = pgerror.Newf(pgcode.Warning, "format %s", wrappedErr.Error())                          // want `err.Error\(\) is passed to pgerror.Newf; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`

	_ = errors.Wrap(anotherErr, wrappedErr.Error())                                          // want `err.Error\(\) is passed to errors.Wrap; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`
	_ = errors.Wrapf(anotherErr, "format %s", wrappedErr.Error())                            // want `err.Error\(\) is passed to errors.Wrapf; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`
	_ = errors.WrapWithDepthf(1, anotherErr, "format %s", wrappedErr.Error())                // want `err.Error\(\) is passed to errors.WrapWithDepthf; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`
	_ = errors.New(wrappedErr.Error())                                                       // want `err.Error\(\) is passed to errors.New; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`
	_ = errors.Newf("format %d %s", 1, wrappedErr.Error())                                   // want `err.Error\(\) is passed to errors.Newf; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`
	_ = errors.NewWithDepthf(1, "format %s", wrappedErr.Error())                             // want `err.Error\(\) is passed to errors.NewWithDepthf; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`
	_ = errors.AssertionFailedf(wrappedErr.Error())                                          // want `err.Error\(\) is passed to errors.AssertionFailedf; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`
	_ = errors.AssertionFailedWithDepthf(1, "format %s", wrappedErr.Error())                 // want `err.Error\(\) is passed to errors.AssertionFailedWithDepthf; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`
	_ = errors.NewAssertionErrorWithWrappedErrf(anotherErr, "format %s", wrappedErr.Error()) // want `err.Error\(\) is passed to errors.NewAssertionErrorWithWrappedErrf; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`

	_ = fmt.Errorf("got %s", wrappedErr)  // want `non-wrapped error is passed to fmt.Errorf; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`
	_ = fmt.Errorf("got %v", wrappedErr)  // want `non-wrapped error is passed to fmt.Errorf; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`
	_ = fmt.Errorf("got %+v", wrappedErr) // want `non-wrapped error is passed to fmt.Errorf; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`
	_ = fmt.Errorf("got %w", wrappedErr)  // this is allowed because of the %w verb`

	_ = pgerror.Wrapf(anotherErr, pgcode.Warning, "format %s", wrappedErr)             // want `non-wrapped error is passed to pgerror.Wrapf; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`
	_ = pgerror.WrapWithDepthf(1, anotherErr, pgcode.Warning, "format %s", wrappedErr) // want `non-wrapped error is passed to pgerror.WrapWithDepthf; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`
	_ = pgerror.Newf(pgcode.Warning, "format %s", wrappedErr)                          // want `non-wrapped error is passed to pgerror.Newf; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`

	_ = errors.Wrapf(anotherErr, "format %v", wrappedErr)                            // want `non-wrapped error is passed to errors.Wrapf; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`
	_ = errors.WrapWithDepthf(1, anotherErr, "format %+v", wrappedErr)               // want `non-wrapped error is passed to errors.WrapWithDepthf; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`
	_ = errors.Newf("format %d %s", 1, wrappedErr)                                   // want `non-wrapped error is passed to errors.Newf; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`
	_ = errors.NewWithDepthf(1, "format %s", wrappedErr)                             // want `non-wrapped error is passed to errors.NewWithDepthf; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`
	_ = errors.AssertionFailedf("format %v", wrappedErr)                             // want `non-wrapped error is passed to errors.AssertionFailedf; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`
	_ = errors.AssertionFailedWithDepthf(1, "format %s", wrappedErr)                 // want `non-wrapped error is passed to errors.AssertionFailedWithDepthf; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`
	_ = errors.NewAssertionErrorWithWrappedErrf(anotherErr, "format %v", wrappedErr) // want `non-wrapped error is passed to errors.NewAssertionErrorWithWrappedErrf; use pgerror.Wrap/errors.Wrap/errors.CombineErrors/errors.WithSecondaryError/errors.NewAssertionErrorWithWrappedErrf instead`

	// nolint:errwrap
	_ = errors.Wrapf(
		wrappedErr,
		"error parsing %s: %v",
		"blah",
		anotherErr,
	)
}
