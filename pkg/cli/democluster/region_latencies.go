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

package democluster

type regionPair struct {
	regionA string
	regionB string
}

var regionToRegionToLatency map[string]map[string]int

func insertPair(pair regionPair, latency int) {
	regionToLatency, ok := regionToRegionToLatency[pair.regionA]
	if !ok {
		regionToLatency = make(map[string]int)
		regionToRegionToLatency[pair.regionA] = regionToLatency
	}
	regionToLatency[pair.regionB] = latency
}

// Round-trip latencies collected from http://cloudping.co on 2019-09-11.
var regionRoundTripLatencies = map[regionPair]int{
	{regionA: "us-east1", regionB: "us-west1"}:     66,
	{regionA: "us-east1", regionB: "europe-west1"}: 64,
	{regionA: "us-west1", regionB: "europe-west1"}: 146,
}

var regionOneWayLatencies = make(map[regionPair]int)

func init() {
	// We record one-way latencies next, because the logic in our delayingConn
	// and delayingListener is in terms of one-way network delays.
	for pair, latency := range regionRoundTripLatencies {
		regionOneWayLatencies[pair] = latency / 2
	}
	regionToRegionToLatency = make(map[string]map[string]int)
	for pair, latency := range regionOneWayLatencies {
		insertPair(pair, latency)
		insertPair(regionPair{
			regionA: pair.regionB,
			regionB: pair.regionA,
		}, latency)
	}
}
