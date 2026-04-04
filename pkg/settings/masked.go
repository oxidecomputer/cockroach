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

package settings

// MaskedSetting is a pseudo-variable constructed on-the-fly by Lookup
// when the actual setting is non-reportable.
type MaskedSetting struct {
	setting NonMaskedSetting
}

var _ Setting = &MaskedSetting{}

// UnderlyingSetting retrieves the actual setting object.
func (s *MaskedSetting) UnderlyingSetting() NonMaskedSetting {
	return s.setting
}

// String hides the underlying value.
func (s *MaskedSetting) String(sv *Values) string {
	// Special case for non-reportable strings: we still want
	// to distinguish empty from non-empty (= customized).
	if st, ok := s.UnderlyingSetting().(*StringSetting); ok && st.String(sv) == "" {
		return ""
	}
	return "<redacted>"
}

// Visibility returns the visibility setting for the underlying setting.
func (s *MaskedSetting) Visibility() Visibility {
	return s.setting.Visibility()
}

// Key returns the key string for the underlying setting.
func (s *MaskedSetting) Key() string {
	return s.setting.Key()
}

// Description returns the description string for the underlying setting.
func (s *MaskedSetting) Description() string {
	return s.setting.Description()
}

// Typ returns the short (1 char) string denoting the type of setting.
func (s *MaskedSetting) Typ() string {
	return s.setting.Typ()
}

// Class returns the class for the underlying setting.
func (s *MaskedSetting) Class() Class {
	return s.setting.Class()
}

// TestingIsReportable is used in testing for reportability.
func TestingIsReportable(s Setting) bool {
	if _, ok := s.(*MaskedSetting); ok {
		return false
	}
	if e, ok := s.(internalSetting); ok {
		return e.isReportable()
	}
	return true
}
