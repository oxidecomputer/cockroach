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

package circuit

import (
	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/redact"
)

// An EventHandler is reported to by circuit breakers.
type EventHandler interface {
	OnTrip(_ *Breaker, prev, cur error)
	OnProbeLaunched(*Breaker)
	OnProbeDone(*Breaker)
	OnReset(*Breaker)
}

// EventLogger is an implementation of EventHandler that relays to a logging
// function. For each event, Log is invoked with a StringBuilder containing
// a log message about the event.
type EventLogger struct {
	Log func(redact.StringBuilder)
}

// OnTrip implements EventHandler. If the previous error is nil, it logs the
// error. If the previous error is not nil and has a different cause than the
// current error, logs a message indicating the old and new error.
func (d *EventLogger) OnTrip(b *Breaker, prev, cur error) {
	var buf redact.StringBuilder
	opts := b.Opts()
	if prev != nil && !errors.Is(errors.Cause(prev), errors.Cause(cur)) {
		buf.Printf("%s: now tripped with error: %s (previously: %s)", opts.Name, cur, prev)
	} else {
		buf.Printf("%s: tripped with error: %s", opts.Name, cur)
	}
	d.Log(buf)
}

// OnProbeLaunched implements EventHandler. It is a no-op.
func (d *EventLogger) OnProbeLaunched(*Breaker) {}

// OnProbeDone implements EventHandler. It is a no-op.
func (d *EventLogger) OnProbeDone(*Breaker) {}

// OnReset implements EventHandler. It logs a message.
func (d *EventLogger) OnReset(b *Breaker) {
	var buf redact.StringBuilder
	buf.Printf("%s: breaker reset", b.Opts().Name)
	d.Log(buf)
}
