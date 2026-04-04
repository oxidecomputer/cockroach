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

package blobs

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
)

// TestBlobServiceClient can be used as a mock BlobClient
// in tests that use nodelocal storage.
func TestBlobServiceClient(externalIODir string) BlobClientFactory {
	return func(ctx context.Context, dialing roachpb.NodeID) (BlobClient, error) {
		return NewLocalClient(externalIODir)
	}
}

// TestEmptyBlobClientFactory can be used as a mock BlobClient
// in tests that create ExternalStorage but do not use
// nodelocal storage.
var TestEmptyBlobClientFactory = func(
	ctx context.Context, dialing roachpb.NodeID,
) (BlobClient, error) {
	return nil, nil
}
