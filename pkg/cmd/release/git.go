// Copyright 2019 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package main

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
)

// findNextVersion returns the next release version for given releaseSeries.
func findNextVersion(releaseSeries string) (string, error) {
	prevReleaseVersion, err := findPreviousRelease(releaseSeries)
	if err != nil {
		return "", fmt.Errorf("cannot find previous release: %w", err)
	}
	nextReleaseVersion, err := bumpVersion(prevReleaseVersion)
	if err != nil {
		return "", fmt.Errorf("cannot bump version: %w", err)
	}
	return nextReleaseVersion, nil
}

func findVersions(text string) []*semver.Version {
	var versions []*semver.Version
	for _, line := range strings.Split(text, "\n") {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}
		// Skip builds before alpha.1
		if strings.Contains(trimmedLine, "-alpha.0000") {
			continue
		}
		version, err := semver.NewVersion(trimmedLine)
		if err != nil {
			fmt.Printf("WARNING: cannot parse version '%s'\n", trimmedLine)
			continue
		}
		versions = append(versions, version)
	}
	return versions
}

// findPreviousRelease finds the latest version tag for a particular release series.
// It ignores non-semantic versions and tags with the alpha.0* suffix.
func findPreviousRelease(releaseSeries string) (string, error) {
	// TODO: filter version using semantic version, not a git pattern
	cmd := exec.Command("git", "tag", "--list", fmt.Sprintf("v%s.*", releaseSeries))
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("cannot get version from tags: %w", err)
	}
	versions := findVersions(string(output))
	if len(versions) == 0 {
		return "", fmt.Errorf("zero versions found")
	}
	sort.Sort(semver.Collection(versions))
	return versions[len(versions)-1].Original(), nil
}

// bumpVersion increases the patch release version (the last digit) of a given version
func bumpVersion(version string) (string, error) {
	semanticVersion, err := semver.NewVersion(version)
	if err != nil {
		return "", fmt.Errorf("cannot parse version: %w", err)
	}
	nextVersion := semanticVersion.IncPatch()
	return nextVersion.Original(), nil
}
