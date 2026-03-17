// Copyright 2017 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package logcrash

import (
	"context"
	"sync/atomic"

	"github.com/cockroachdb/cockroach/pkg/build"
	"github.com/cockroachdb/cockroach/pkg/settings"
	"github.com/cockroachdb/cockroach/pkg/util/envutil"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/log/severity"
	"github.com/cockroachdb/errors"
)

// The call stack here is usually:
// - ReportPanic
// - RecoverAndReport
// - panic()
// so ReportPanic should pop three frames.
const depthForRecoverAndReportPanic = 3

var (
	// PanicOnAssertions wraps "debug.panic_on_failed_assertions"
	PanicOnAssertions = settings.RegisterBoolSetting(
		settings.TenantWritable,
		"debug.panic_on_failed_assertions",
		"panic when an assertion fails rather than reporting",
		false,
	)

	// ReportSensitiveDetails enables reporting of unanonymized data.
	//
	// This should not be used by anyone unwilling to share the whole cluster
	// data with Cockroach Labs and various cloud services.
	ReportSensitiveDetails = envutil.EnvOrDefaultBool("COCKROACH_REPORT_SENSITIVE_DETAILS", false)

	// globalSettings stores a global reference to a *setting.Values container;
	// used for code paths where the container is not available.
	globalSettings atomic.Value
)

// SetGlobalSettings sets the *settings.Values container that will be refreshed
// at runtime -- ideally we should have no other *Values containers floating
// around, as they will be stale / lies.
func SetGlobalSettings(v *settings.Values) {
	globalSettings.Store(v)
}

func getGlobalSettings() *settings.Values {
	if ptr := globalSettings.Load(); ptr != nil {
		return ptr.(*settings.Values)
	}
	return nil
}

// ReportPanicWithGlobalSettings is a variant of ReportPanic that uses the
// *settings.Values that was set using SetGlobalSettings(). Does nothing if that
// function was not called.
//
// Should be used only when strictly necessary; use ReportPanic whenever we have
// access to the settings.
func ReportPanicWithGlobalSettings(ctx context.Context, r interface{}, depth int) {
	if sv := getGlobalSettings(); sv != nil {
		ReportPanic(ctx, sv, r, depth+1)
	}
}

// RecoverAndReportPanic can be invoked on goroutines that run with
// stderr redirected to logs to ensure the user gets informed on the
// real stderr a panic has occurred.
func RecoverAndReportPanic(ctx context.Context, sv *settings.Values) {
	if r := recover(); r != nil {
		ReportPanic(ctx, sv, r, depthForRecoverAndReportPanic)
		panic(r)
	}
}

// RecoverAndReportNonfatalPanic is an alternative RecoverAndReportPanic that
// does not re-panic in Release builds.
func RecoverAndReportNonfatalPanic(ctx context.Context, sv *settings.Values) {
	if r := recover(); r != nil {
		ReportPanic(ctx, sv, r, depthForRecoverAndReportPanic)
		if !build.IsRelease() || PanicOnAssertions.Get(sv) {
			panic(r)
		}
	}
}

// ReportPanic reports a panic has occurred on the real stderr.
//
// Note that ReportPanic() does not assume that the panic object
// will be left uncaught to terminate the process. For example,
// at the time of this writing, ReportPanic() is called from
// RecoverAndReportNonfatalPanic() above.
func ReportPanic(ctx context.Context, sv *settings.Values, r interface{}, depth int) {
	// Announce the panic has occurred to all places. The purpose
	// of this call is threefold:
	// - it ensures there's a notice on the terminal, in case
	//   logging would only go to file otherwise;
	// - it ensures there's a notice on the output file, in
	//   case the panic is uncaught and internal stderr
	//   writes by the Go runtime have not been set up to
	//   redirect to a separate log file.
	// - it places a timestamp in front of the panic object,
	//   in case the various configuration options make
	//   the Go runtime solely responsible for printing
	//   out the panic object to the log file.
	//   (The go runtime doesn't time stamp its output.)
	//
	// Note that this code will cause the panic object to be printed
	// twice in some cases (specifically, when the panic is uncaught).
	// A previous version of this code was trying to prevent the
	// double print and was failing to do so effectively, causing
	// instead panics to be lost in the case where they were
	// recovered (eg via RecoverAndReportNonfatalPanic()).
	//
	// To properly prevent double prints, the API could be changed to
	// indicate whether the Go runtime will *eventually* print the panic
	// on its own. Unfortunately, this is a bit hard to do, as the
	// caller of ReportPanic() may not be in a position to know for
	// sure, whether some other caller further in the call stack is
	// catching the panic object in the end or not.
	panicErr := PanicAsError(depth+1, r)
	log.Ops.Shoutf(ctx, severity.ERROR, "a panic has occurred!\n%+v", panicErr)

	// Ensure that the logs are flushed before letting a panic
	// terminate the server.
	log.Flush()
}

// PanicAsError turns r into an error if it is not one already.
func PanicAsError(depth int, r interface{}) error {
	if err, ok := r.(error); ok {
		return errors.WithStackDepth(err, depth+1)
	}
	return errors.NewWithDepthf(depth+1, "panic: %v", r)
}

// ReportOrPanic either reports an error to sentry, if run from a release
// binary, or panics, if triggered in tests. This is intended to be used for
// failing assertions which are recoverable but serious enough to report and to
// cause tests to fail.
//
// Like SendCrashReport, the format string should not contain any sensitive
// data, and unsafe reportables will be redacted before reporting.
func ReportOrPanic(
	ctx context.Context, sv *settings.Values, format string, reportables ...interface{},
) {
	err := errors.Newf(format, reportables...)
	if !build.IsRelease() || (sv != nil && PanicOnAssertions.Get(sv)) {
		panic(err)
	}
	log.Warningf(ctx, "%v", err)
}
