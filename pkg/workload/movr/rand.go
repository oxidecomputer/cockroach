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

package movr

import (
	"encoding/json"

	"golang.org/x/exp/rand"
)

const numerals = `1234567890`

var vehicleTypes = [...]string{`skateboard`, `bike`, `scooter`}
var vehicleColors = [...]string{`red`, `yellow`, `blue`, `green`, `black`}
var bikeBrands = [...]string{
	`Merida`, `Fuji`, `Cervelo`, `Pinarello`, `Santa Cruz`, `Kona`, `Schwinn`}

func randString(rng *rand.Rand, length int, alphabet string) string {
	buf := make([]byte, length)
	for i := range buf {
		buf[i] = alphabet[rng.Intn(len(alphabet))]
	}
	return string(buf)
}

func randCreditCard(rng *rand.Rand) string {
	return randString(rng, 10, numerals)
}

func randVehicleType(rng *rand.Rand) string {
	return vehicleTypes[rng.Intn(len(vehicleTypes))]
}

func randVehicleStatus(rng *rand.Rand) string {
	r := rng.Intn(100)
	switch {
	case r < 40:
		return `available`
	case r < 95:
		return `in_use`
	default:
		return `lost`
	}
}

func randLatLong(rng *rand.Rand) (float64, float64) {
	lat, long := float64(-180+rng.Intn(360)), float64(-90+rng.Intn(180))
	return lat, long
}

func randCity(rng *rand.Rand) string {
	idx := rng.Int31n(int32(len(cities)))
	return cities[idx].city
}

func randVehicleMetadata(rng *rand.Rand, vehicleType string) string {
	m := map[string]string{
		`color`: vehicleColors[rng.Intn(len(vehicleColors))],
	}
	switch vehicleType {
	case `bike`:
		m[`brand`] = bikeBrands[rng.Intn(len(bikeBrands))]
	}
	j, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(j)
}
