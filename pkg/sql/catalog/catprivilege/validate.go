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

package catprivilege

import (
	"strings"

	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/catpb"
	"github.com/cockroachdb/cockroach/pkg/sql/privilege"
)

// Validate validates a privilege descriptor.
func Validate(
	p catpb.PrivilegeDescriptor, objectNameKey catalog.NameKey, objectType privilege.ObjectType,
) error {
	return p.Validate(
		objectNameKey.GetParentID(),
		objectType,
		objectNameKey.GetName(),
		allowedSuperuserPrivileges(objectNameKey),
	)
}

// ValidateSuperuserPrivileges validates superuser privileges.
func ValidateSuperuserPrivileges(
	p catpb.PrivilegeDescriptor, objectNameKey catalog.NameKey, objectType privilege.ObjectType,
) error {
	return p.ValidateSuperuserPrivileges(
		objectNameKey.GetParentID(),
		objectType,
		objectNameKey.GetName(),
		allowedSuperuserPrivileges(objectNameKey),
	)
}

// ValidateDefaultPrivileges validates default privileges.
func ValidateDefaultPrivileges(p catpb.DefaultPrivilegeDescriptor) error {
	return p.Validate()
}

func allowedSuperuserPrivileges(objectNameKey catalog.NameKey) privilege.List {
	privs := SystemSuperuserPrivileges(objectNameKey)
	if privs != nil {
		return privs
	}
	// Cluster restores move certain system tables to a higher ID to prevent
	// conflicts with non-system descriptors that are going to be restored. The
	// newly created tables in the system database will be given ReadWrite
	// privileges.
	//
	// TODO(adityamaru,dt): Remove once we fix the handling of dynamic system
	// table IDs during restore.
	if objectNameKey.GetParentID() == keys.SystemDatabaseID &&
		strings.Contains(objectNameKey.GetName(), RestoreCopySystemTablePrefix) {
		return privilege.ReadWriteData
	}
	return catpb.DefaultSuperuserPrivileges
}
