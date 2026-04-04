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

package storage

import (
	"io/ioutil"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/storage/fs"
	"github.com/cockroachdb/cockroach/pkg/util/protoutil"
	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/errors/oserror"
	"github.com/cockroachdb/pebble/vfs"
)

// MinVersionFilename is the name of the file containing a marshaled
// roachpb.Version that can be updated during storage-related migrations
// and checked on startup to determine if we can safely use a
// backwards-incompatible feature.
const MinVersionFilename = "STORAGE_MIN_VERSION"

// writeMinVersionFile writes the provided version to disk. The caller must
// guarantee that the version will never be downgraded below the given version.
func writeMinVersionFile(atomicRenameFS vfs.FS, dir string, version roachpb.Version) error {
	// TODO(jackson): Assert that atomicRenameFS supports atomic renames
	// once Pebble is bumped to the appropriate SHA.
	if version == (roachpb.Version{}) {
		return errors.New("min version should not be empty")
	}
	ok, err := MinVersionIsAtLeastTargetVersion(atomicRenameFS, dir, version)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	b, err := protoutil.Marshal(&version)
	if err != nil {
		return err
	}
	filename := atomicRenameFS.PathJoin(dir, MinVersionFilename)
	if err := fs.SafeWriteToFile(atomicRenameFS, dir, filename, b); err != nil {
		return err
	}
	return nil
}

// MinVersionIsAtLeastTargetVersion returns whether the min version recorded
// on disk is at least the target version.
func MinVersionIsAtLeastTargetVersion(
	atomicRenameFS vfs.FS, dir string, target roachpb.Version,
) (bool, error) {
	// TODO(jackson): Assert that atomicRenameFS supports atomic renames
	// once Pebble is bumped to the appropriate SHA.
	if target == (roachpb.Version{}) {
		return false, errors.New("target version should not be empty")
	}
	minVersion, err := getMinVersion(atomicRenameFS, dir)
	if err != nil {
		return false, err
	}
	if minVersion == (roachpb.Version{}) {
		return false, nil
	}
	return !minVersion.Less(target), nil
}

// getMinVersion returns the min version recorded on disk if the min version
// file exists and nil otherwise.
func getMinVersion(atomicRenameFS vfs.FS, dir string) (roachpb.Version, error) {
	// TODO(jackson): Assert that atomicRenameFS supports atomic renames
	// once Pebble is bumped to the appropriate SHA.

	filename := atomicRenameFS.PathJoin(dir, MinVersionFilename)
	f, err := atomicRenameFS.Open(filename)
	if oserror.IsNotExist(err) {
		return roachpb.Version{}, nil
	}
	if err != nil {
		return roachpb.Version{}, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return roachpb.Version{}, err
	}
	version := roachpb.Version{}
	if err := protoutil.Unmarshal(b, &version); err != nil {
		return version, err
	}
	return version, nil
}
