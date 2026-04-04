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

/*
Package impl is a stub package that imports all of the concrete
implementations of the various cloud storage providers to trigger their
initialization-time registration with the cloud storage provider registry.
*/
package impl

import (
	// Import all the cloud provider packages to register them.
	_ "github.com/cockroachdb/cockroach/pkg/cloud/amazon"
	_ "github.com/cockroachdb/cockroach/pkg/cloud/azure"
	_ "github.com/cockroachdb/cockroach/pkg/cloud/gcp"
	_ "github.com/cockroachdb/cockroach/pkg/cloud/httpsink"
	_ "github.com/cockroachdb/cockroach/pkg/cloud/nodelocal"
	_ "github.com/cockroachdb/cockroach/pkg/cloud/nullsink"
	_ "github.com/cockroachdb/cockroach/pkg/cloud/userfile"
)
