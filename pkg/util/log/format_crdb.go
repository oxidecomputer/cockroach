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

package log

import (
	"hash/adler32"
	"io"

	"github.com/cockroachdb/cockroach/pkg/util/log/logpb"
	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/ttycolor"
)

const severityChar = "IWEF"

// MessageTimeFormat is the format of the timestamp in log message headers of crdb formatted logs.
// as used in time.Parse and time.Format.
const MessageTimeFormat = "060102 15:04:05.999999"

// FormatLegacyEntry writes the contents of the legacy log entry struct to the specified writer.
func FormatLegacyEntry(e logpb.Entry, w io.Writer) error {
	return FormatLegacyEntryWithOptionalColors(e, w, nil /* cp */)
}

// FormatLegacyEntryWithOptionalColors is like FormatLegacyEntry but the caller can specify
// a color profile.
func FormatLegacyEntryWithOptionalColors(e logpb.Entry, w io.Writer, cp ttycolor.Profile) error {
	buf := formatLogEntryInternalV1(e, false /* isHeader */, true /* showCounter */, cp)
	defer putBuffer(buf)
	_, err := w.Write(buf.Bytes())
	return err
}

// FormatLegacyEntryPrefix writes a color-decorated prefix to the specified
// writer. The color is rendered in the background of the prefix and is chosen
// from an arbitrary but deterministic mapping from the prefix bytes to the
// color profile entries.
func FormatLegacyEntryPrefix(prefix []byte, w io.Writer, cp ttycolor.Profile) (err error) {
	if prefix == nil {
		return nil
	}

	if cp != nil {
		code := ttycolor.PickArbitraryColor(adler32.Checksum(prefix))
		if _, err = w.Write(cp.BackgroundColorSequence(code)); err != nil {
			return err
		}
		defer func() {
			_, errReset := w.Write(cp[ttycolor.Reset])
			err = errors.CombineErrors(err, errReset)
		}()
	}

	_, err = w.Write(prefix)
	return err
}
