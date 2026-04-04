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

package descpb

type Descriptor struct{}

type TableDescriptor struct{}

type DatabaseDescriptor struct{}

type TypeDescriptor struct{}

type SchemaDescriptor struct{}

func (m *Descriptor) GetTable() *TableDescriptor {
	return nil
}

func (m *Descriptor) GetDatabase() *DatabaseDescriptor {
	return nil
}

func (m *Descriptor) GetType() *TypeDescriptor {
	return nil
}

func (m *Descriptor) GetSchema() *SchemaDescriptor {
	return nil
}
