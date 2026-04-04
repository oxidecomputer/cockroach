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

package tests

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
)

var pythonUnitTestOutputRegex = regexp.MustCompile(`(?P<name>.*) \((?P<class>.*)\) \.\.\. (?P<result>[^'"]*?)(?: u?['"](?P<reason>.*)['"])?$`)

func (r *ormTestsResults) parsePythonUnitTestOutput(
	input []byte, expectedFailures blocklist, ignoredList blocklist,
) {
	scanner := bufio.NewScanner(bytes.NewReader(input))
	for scanner.Scan() {
		match := pythonUnitTestOutputRegex.FindStringSubmatch(scanner.Text())
		if match != nil {
			groups := map[string]string{}
			for i, name := range match {
				groups[pythonUnitTestOutputRegex.SubexpNames()[i]] = name
			}
			test := fmt.Sprintf("%s.%s", groups["class"], groups["name"])
			skipped := groups["result"] == "skipped" || groups["result"] == "expected failure"
			skipReason := ""
			if skipped {
				skipReason = groups["reason"]
			}
			pass := groups["result"] == "ok" || groups["result"] == "unexpected success"
			r.allTests = append(r.allTests, test)

			ignoredIssue, expectedIgnored := ignoredList[test]
			issue, expectedFailure := expectedFailures[test]
			switch {
			case expectedIgnored:
				r.results[test] = fmt.Sprintf("--- SKIP: %s due to %s (expected)", test, ignoredIssue)
				r.ignoredCount++
			case skipped && expectedFailure:
				r.results[test] = fmt.Sprintf("--- SKIP: %s due to %s (unexpected)", test, skipReason)
				r.unexpectedSkipCount++
			case skipped:
				r.results[test] = fmt.Sprintf("--- SKIP: %s due to %s (expected)", test, skipReason)
				r.skipCount++
			case pass && !expectedFailure:
				r.results[test] = fmt.Sprintf("--- PASS: %s (expected)", test)
				r.passExpectedCount++
			case pass && expectedFailure:
				r.results[test] = fmt.Sprintf("--- PASS: %s - %s (unexpected)",
					test, maybeAddGithubLink(issue),
				)
				r.passUnexpectedCount++
			case !pass && expectedFailure:
				r.results[test] = fmt.Sprintf("--- FAIL: %s - %s (expected)",
					test, maybeAddGithubLink(issue),
				)
				r.failExpectedCount++
				r.currentFailures = append(r.currentFailures, test)
			case !pass && !expectedFailure:
				r.results[test] = fmt.Sprintf("--- FAIL: %s (unexpected)", test)
				r.failUnexpectedCount++
				r.currentFailures = append(r.currentFailures, test)
			}
			r.runTests[test] = struct{}{}
		}
	}
}
