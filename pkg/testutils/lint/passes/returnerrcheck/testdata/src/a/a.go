// Copyright 2019 The Cockroach Authors.
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
	"errors"
	"fmt"
	"log"
)

func noReturn() {
	err := errors.New("foo")
	if err != nil {
		return
	}
	_ = func() error {
		if err != nil {
			return nil // want `unexpected nil error return after checking for a non-nil error`
		}
		return nil
	}()
	return
}

func SingleReturn() error {
	err := errors.New("foo")
	if err != nil {
		return nil // want `unexpected nil error return after checking for a non-nil error`
	}
	if err != nil || false {
		return nil // want `unexpected nil error return after checking for a non-nil error`
	}
	if false || err != nil {
		return nil // want `unexpected nil error return after checking for a non-nil error`
	}
	if err != nil && false {
		return nil // want `unexpected nil error return after checking for a non-nil error`
	}
	if false && err != nil {
		return nil // want `unexpected nil error return after checking for a non-nil error`
	}
	if false || err != nil && false {
		return nil // want `unexpected nil error return after checking for a non-nil error`
	}
	if true && false && err != nil {
		return nil // want `unexpected nil error return after checking for a non-nil error`
	}
	if err != nil {
		//nolint:returnerrcheck
		return nil
	}
	if err != nil {
		// nolint:returnerrcheck
		return nil
	}
	if err != nil {
		return nil // nolint:returnerrcheck
	}

	// The below nolint comments is not close enough to the return statement.

	// nolint:returnerrcheck
	if err != nil {
		return nil // want `unexpected nil error return after checking for a non-nil error`
	}

	if err != nil {
		// nolint:returnerrcheck
		fmt.Println("hmm")

		return nil // want `unexpected nil error return after checking for a non-nil error`
	}
	if err != nil {
		return nil //nolint:returnerrcheck
	}
	return nil
}

func MultipleReturns() (int, error) {
	err := errors.New("foo")
	if err != nil {
		if true {
			return 0, nil
		}
		return 0, nil // want `unexpected nil error return after checking for a non-nil error`
	}
	if err != nil {
		//nolint:returnerrcheck
		return 0, nil
	}
	return -1, nil
}

func AcceptableBehavior() {
	type structThing struct {
		err error
	}
	_ = func() error {
		var t structThing
		// Intentionally don't error if the check is not for an ident directly but
		// rather is for a field.
		if t.err != nil {
			return nil
		}
		return nil
	}
	_ = func() error {
		if err := errors.New("foo"); err != nil {
			log.Printf("%v", err)
			return nil
		}
		return nil
	}
	var t structThing
	_ = func() error {
		if err := errors.New("foo"); err != nil {
			t.err = err
			return nil
		}
		return nil
	}
	func() error {
		if err := errors.New("foo"); err != nil {
			t = structThing{err: err}
			return nil
		}
		return nil
	}()
	func() error {
		if err := errors.New("foo"); err != nil || true {
			if err != nil {
				return err
			}
			return nil
		}
		return nil
	}()
	func() (string, error) {
		if err := errors.New("foo"); err != nil {
			return err.Error(), nil
		}
		return "", nil
	}()
}
