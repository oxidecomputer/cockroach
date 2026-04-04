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

package timeutil

// FullTimeFormat is the time format used to display any unknown timestamp
// type, and always shows the full time zone offset.
const FullTimeFormat = "2006-01-02 15:04:05.999999-07:00:00"

// TimestampWithTZFormat is the time format used to display
// timestamps with a time zone offset. The minutes and seconds
// offsets are only added if they are non-zero.
const TimestampWithTZFormat = "2006-01-02 15:04:05.999999-07"

// TimestampWithoutTZFormat is the time format used to display
// timestamps without a time zone offset. The minutes and seconds
// offsets are only added if they are non-zero.
const TimestampWithoutTZFormat = "2006-01-02 15:04:05.999999"

// TimeWithTZFormat is the time format used to display a time
// with a time zone offset.
const TimeWithTZFormat = "15:04:05.999999-07"

// TimeWithoutTZFormat is the time format used to display a time
// without a time zone offset.
const TimeWithoutTZFormat = "15:04:05.999999"

// DateFormat is the time format used to display a date.
const DateFormat = "2006-01-02"
