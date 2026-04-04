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

package spec

import "fmt"

// AWSMachineType selects a machine type given the desired number of CPUs.
func AWSMachineType(cpus int, highmem bool) string {
	// TODO(erikgrinaker): These have significantly less RAM than
	// their GCE counterparts. Consider harmonizing them.
	family := "c5d" // 2 GB RAM per CPU
	if highmem {
		family = "m5d" // 4 GB RAM per CPU
	}

	var size string
	switch {
	case cpus <= 2:
		size = "large"
	case cpus <= 4:
		size = "xlarge"
	case cpus <= 8:
		size = "2xlarge"
	case cpus <= 16:
		size = "4xlarge"
	case cpus <= 36:
		size = "9xlarge"
	case cpus <= 72:
		size = "18xlarge"
	case cpus <= 96:
		size = "24xlarge"
	default:
		panic(fmt.Sprintf("no aws machine type with %d cpus", cpus))
	}

	// There is no c5d.24xlarge.
	if family == "c5d" && size == "24xlarge" {
		family = "m5d"
	}

	return fmt.Sprintf("%s.%s", family, size)
}

// GCEMachineType selects a machine type given the desired number of CPUs.
func GCEMachineType(cpus int, highmem bool) string {
	// TODO(peter): This is awkward: at or below 16 cpus, use n1-standard so that
	// the machines have a decent amount of RAM. We could use custom machine
	// configurations, but the rules for the amount of RAM per CPU need to be
	// determined (you can't request any arbitrary amount of RAM).
	series := "n1"
	kind := "standard" // 3.75 GB RAM per CPU
	if highmem {
		kind = "highmem" // 6.5 GB RAM per CPU
	} else if cpus > 16 {
		kind = "highcpu" // 0.9 GB RAM per CPU
	}
	return fmt.Sprintf("%s-%s-%d", series, kind, cpus)
}

// AzureMachineType selects a machine type given the desired number of CPUs.
func AzureMachineType(cpus int, highmem bool) string {
	if highmem {
		panic("highmem not implemented for Azure")
	}
	switch {
	case cpus <= 2:
		return "Standard_D2_v3"
	case cpus <= 4:
		return "Standard_D4_v3"
	case cpus <= 8:
		return "Standard_D8_v3"
	case cpus <= 16:
		return "Standard_D16_v3"
	case cpus <= 36:
		return "Standard_D32_v3"
	case cpus <= 48:
		return "Standard_D48_v3"
	case cpus <= 64:
		return "Standard_D64_v3"
	default:
		panic(fmt.Sprintf("no azure machine type with %d cpus", cpus))
	}
}
