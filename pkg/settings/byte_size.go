// Copyright 2017 The Cockroach Authors.
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

package settings

import (
	"github.com/cockroachdb/cockroach/pkg/util/humanizeutil"
	"github.com/cockroachdb/errors"
)

// ByteSizeSetting is the interface of a setting variable that will be
// updated automatically when the corresponding cluster-wide setting
// of type "bytesize" is updated.
type ByteSizeSetting struct {
	IntSetting
}

var _ numericSetting = &ByteSizeSetting{}

// Typ returns the short (1 char) string denoting the type of setting.
func (*ByteSizeSetting) Typ() string {
	return "z"
}

func (b *ByteSizeSetting) String(sv *Values) string {
	return string(humanizeutil.IBytes(b.Get(sv)))
}

// DecodeToString decodes and renders an encoded value.
func (b *ByteSizeSetting) DecodeToString(encoded string) (string, error) {
	iv, err := b.DecodeValue(encoded)
	if err != nil {
		return "", err
	}
	return string(humanizeutil.IBytes(iv)), nil
}

// WithPublic sets public visibility and can be chained.
func (b *ByteSizeSetting) WithPublic() *ByteSizeSetting {
	b.SetVisibility(Public)
	return b
}

// RegisterByteSizeSetting defines a new setting with type bytesize and any
// supplied validation function(s).
func RegisterByteSizeSetting(
	class Class, key, desc string, defaultValue int64, validateFns ...func(int64) error,
) *ByteSizeSetting {

	var validateFn func(int64) error
	if len(validateFns) > 0 {
		validateFn = func(v int64) error {
			for _, fn := range validateFns {
				if err := fn(v); err != nil {
					return errors.Wrapf(err, "invalid value for %s", key)
				}
			}
			return nil
		}
	}

	if validateFn != nil {
		if err := validateFn(defaultValue); err != nil {
			panic(errors.Wrap(err, "invalid default"))
		}
	}
	setting := &ByteSizeSetting{IntSetting{
		defaultValue: defaultValue,
		validateFn:   validateFn,
	}}
	register(class, key, desc, setting)
	return setting
}
