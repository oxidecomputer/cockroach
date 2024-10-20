// Copyright 2015 The Cockroach Authors.
// Copyright 2024 Oxide Computer Company
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package caller

import (
	"fmt"
	"path"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/cockroachdb/cockroach/pkg/util/syncutil"
)

type cachedLookup struct {
	file string
	line int
	fun  string
}

var dummyLookup = cachedLookup{file: "???", line: 1, fun: "???"}

// A CallResolver is a helping hand around runtime.Caller() to look up file,
// line and name of the calling function. CallResolver caches the results of
// its lookups and strips the uninteresting prefix from both the caller's
// location and name; see NewCallResolver().
type CallResolver struct {
	mu    syncutil.Mutex
	cache map[uintptr]*cachedLookup
	re    *regexp.Regexp
}

var reStripNothing = regexp.MustCompile(`^$`)

// findPackageRoot finds the root of the `github.com/cockroachdb/cockraoch`
// package. (This is equivalent to the root of the Git repository.)
//
// For example:
//
//     /home/kena/src/go/src/github.com/cockroachdb/cockroach/pkg/util/caller
//     ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ package root
//
// The first return value is false if the paths could not be determined.
func findPackageRoot() (ok bool, pkgRoot string) {
	pcs := make([]uintptr, 1)
	if runtime.Callers(1, pcs[:]) < 1 {
		return false, ""
	}
	frame, _ := runtime.CallersFrames(pcs).Next()

	// frame.Function is the name of the function prefixed by its
	// *symbolic* package path.
	// For example:
	//     github.com/cockroachdb/cockroach/pkg/util/caller.findFileAndPackageRoot
	funcName := frame.Function

	modulePath := strings.TrimSuffix(funcName, ".findPackageRoot")
	innerPath := strings.TrimPrefix(modulePath, "github.com/cockroachdb/cockroach")
	if innerPath == modulePath {
		// `modulePath` does not start with "github.com/cockroachdb/cockroach".
		return false, ""
	}

	// frame.File is the name of the file on the filesystem.
	// For example:
	//   /home/kena/src/go/src/github.com/cockroachdb/cockroach/pkg/util/caller/resolver.go
	//
	// (or, in a Bazel sandbox)
	//   github.com/cockroachdb/cockroach/pkg/util/caller/resolver.go
	//
	// srcRoot is its immediate parent directory.
	srcRoot := path.Dir(frame.File)

	// Coverage tests report back as `[...]/util/caller/_test/_obj_test`;
	// strip back to this package's directory.
	if !strings.HasSuffix(srcRoot, "/caller") {
		// This trims the last component.
		srcRoot = path.Dir(srcRoot)
	}

	if !strings.HasSuffix(srcRoot, "/caller") {
		// If we are not finding the current package in the path, this is
		// indicative of a bug in this code; either:
		//
		// - the name of the function was changed without updating the TrimSuffix
		//   call above.
		// - the package was renamed without updating the two HasSuffix calls
		//   above.
		panic(fmt.Sprintf("cannot find self package: expected .../caller, got %q", srcRoot))
	}

	pkgRoot = strings.TrimSuffix(srcRoot, innerPath)
	if pkgRoot == srcRoot {
		return false, ""
	} else {
		return true, pkgRoot
	}
}

// defaultRE strips as follows:
//
// - <pkgroot>/(pkg/)?module/submodule/file.go
//   -> module/submodule/file.go
//
// - <pkgroot>/vendor/<otherpkg>/path/to/file
//   -> vendor/<otherpkg>/path/to/file
//
// It falls back to stripping nothing when it's unable to look up its
// own location via runtime.Caller().
var defaultRE = func() *regexp.Regexp {
	ok, pkgRoot := findPackageRoot()
	if !ok {
		return reStripNothing
	}

	pkgStrip := regexp.QuoteMeta(pkgRoot) + "/(?:pkg/)?(.*)"
	return regexp.MustCompile(pkgStrip)
}()

var defaultCallResolver = NewCallResolver(defaultRE)

// Lookup returns the (reduced) file, line and function of the caller at the
// requested depth, using a default call resolver which drops the path of
// the project repository.
func Lookup(depth int) (file string, line int, fun string) {
	return defaultCallResolver.Lookup(depth + 1)
}

// NewCallResolver returns a CallResolver. The supplied pattern must specify a
// valid regular expression and is used to format the paths returned by
// Lookup(): If submatches are specified, their concatenation forms the path,
// otherwise the match of the whole expression is used. Paths which do not
// match at all are left unchanged.
// TODO(bdarnell): don't strip paths at lookup time, but at display time;
// need better handling for callers such as x/tools/something.
func NewCallResolver(re *regexp.Regexp) *CallResolver {
	return &CallResolver{
		cache: map[uintptr]*cachedLookup{},
		re:    re,
	}
}

var uintptrSlPool = sync.Pool{
	New: func() interface{} {
		sl := make([]uintptr, 1)
		return &sl
	},
}

// Lookup returns the (reduced) file, line and function of the caller at the
// requested depth.
func (cr *CallResolver) Lookup(depth int) (_file string, _line int, _fun string) {
	sl := uintptrSlPool.Get().(*[]uintptr)
	// NB: +2 for Callers, +1 for Caller (historical reasons)
	ok := runtime.Callers(depth+2, *sl) == 1
	pc := (*sl)[0]
	uintptrSlPool.Put(sl)
	sl = nil // prevent reuse
	if !ok || cr == nil {
		return dummyLookup.file, dummyLookup.line, dummyLookup.fun
	}

	cr.mu.Lock()
	defer cr.mu.Unlock()
	if v, okCache := cr.cache[pc]; okCache {
		return v.file, v.line, v.fun
	}
	// Now do the expensive thing which we intend to cache.
	frame, _ := runtime.CallersFrames([]uintptr{pc}).Next()
	if matches := cr.re.FindStringSubmatch(frame.File); matches != nil {
		if len(matches) == 1 {
			frame.File = matches[0]
		} else {
			// NB: "path" is used here (and elsewhere in this file) over
			// "path/filepath" because runtime.Caller always returns unix paths.
			frame.File = path.Join(matches[1:]...)
		}
	}

	if indDot := strings.LastIndexByte(frame.Function, '.'); indDot != -1 {
		frame.Function = frame.Function[indDot+1:]
	}
	cr.cache[pc] = &cachedLookup{file: frame.File, line: frame.Line, fun: frame.Function}
	return frame.File, frame.Line, frame.Function
}
