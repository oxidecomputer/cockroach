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

package descmarshaltest

import "github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"

func F() {
	var d descpb.Descriptor
	d.GetDatabase() // want `Illegal call to Descriptor.GetDatabase\(\), see descpb.FromDescriptorWithMVCCTimestamp\(\)`

	//nolint:descriptormarshal
	d.GetDatabase()

	//nolint:descriptormarshal
	d.GetDatabase()

	d.GetTable() // want `Illegal call to Descriptor.GetTable\(\), see descpb.FromDescriptorWithMVCCTimestamp\(\)`

	//nolint:descriptormarshal
	d.GetTable()

	//nolint:descriptormarshal
	d.GetTable()

	d.GetType() // want `Illegal call to Descriptor.GetType\(\), see descpb.FromDescriptorWithMVCCTimestamp\(\)`

	//nolint:descriptormarshal
	d.GetType()

	//nolint:descriptormarshal
	d.GetType()

	d.GetSchema() // want `Illegal call to Descriptor.GetSchema\(\), see descpb.FromDescriptorWithMVCCTimestamp\(\)`

	//nolint:descriptormarshal
	d.GetSchema()

	//nolint:descriptormarshal
	d.GetSchema()

	// nolint:descriptormarshal
	if t := d.GetTable(); t != nil {
		panic("foo")
	}

	if t := d.
		// nolint:descriptormarshal
		GetTable(); t != nil {
		panic("foo")
	}

	if t :=
		// nolint:descriptormarshal
		d.GetTable(); t != nil {
		panic("foo")
	}

	if t := d.GetTable(); t != // want `Illegal call to Descriptor.GetTable\(\), see descpb.FromDescriptorWithMVCCTimestamp\(\)`
		// nolint:descriptormarshal
		nil {
		panic("foo")
	}

	// It does not work to put the comment as an inline with the preamble to an
	// if statement.
	if t := d.GetTable(); t != nil { // nolint:descriptormarshal // want `Illegal call to Descriptor.GetTable\(\), see descpb.FromDescriptorWithMVCCTimestamp\(\)`
		panic("foo")
	}

	if t := d.GetTable(); t != nil { // want `Illegal call to Descriptor.GetTable\(\), see descpb.FromDescriptorWithMVCCTimestamp\(\)`
		panic("foo")
	}
}
