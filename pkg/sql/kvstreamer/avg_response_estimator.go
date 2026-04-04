// Copyright 2022 The Cockroach Authors.
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

package kvstreamer

// avgResponseEstimator is a helper that estimates the average size of responses
// received by the Streamer. It is **not** thread-safe.
type avgResponseEstimator struct {
	// responseBytes tracks the total footprint of all responses that the
	// Streamer has already received.
	responseBytes int64
	numResponses  int64
}

// TODO(yuzefovich): use the optimizer-driven estimates.
const initialAvgResponseSize = 1 << 10 // 1KiB

func (e *avgResponseEstimator) getAvgResponseSize() int64 {
	if e.numResponses == 0 {
		return initialAvgResponseSize
	}
	// TODO(yuzefovich): we currently use a simple average over the received
	// responses, but it is likely to be suboptimal because it would be unfair
	// to "large" batches that come in late (i.e. it would not be reactive
	// enough). Consider using another function here.
	return e.responseBytes / e.numResponses
}

// update updates the actual information of the estimator based on numResponses
// responses that took up responseBytes bytes and correspond to a single
// BatchResponse.
func (e *avgResponseEstimator) update(responseBytes int64, numResponses int64) {
	e.responseBytes += responseBytes
	e.numResponses += numResponses
}
