// Copyright 2020 The Cockroach Authors.
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

package tree

// AlterDatabaseOwner represents a ALTER DATABASE OWNER TO statement.
type AlterDatabaseOwner struct {
	Name  Name
	Owner RoleSpec
}

// Format implements the NodeFormatter interface.
func (node *AlterDatabaseOwner) Format(ctx *FmtCtx) {
	ctx.WriteString("ALTER DATABASE ")
	ctx.FormatNode(&node.Name)
	ctx.WriteString(" OWNER TO ")
	ctx.FormatNode(&node.Owner)
}

// AlterDatabaseAddRegion represents a ALTER DATABASE ADD REGION statement.
type AlterDatabaseAddRegion struct {
	Name        Name
	Region      Name
	IfNotExists bool
}

var _ Statement = &AlterDatabaseAddRegion{}

// Format implements the NodeFormatter interface.
func (node *AlterDatabaseAddRegion) Format(ctx *FmtCtx) {
	ctx.WriteString("ALTER DATABASE ")
	ctx.FormatNode(&node.Name)
	ctx.WriteString(" ADD REGION ")
	if node.IfNotExists {
		ctx.WriteString("IF NOT EXISTS ")
	}
	ctx.FormatNode(&node.Region)
}

// AlterDatabaseDropRegion represents a ALTER DATABASE DROP REGION statement.
type AlterDatabaseDropRegion struct {
	Name     Name
	Region   Name
	IfExists bool
}

var _ Statement = &AlterDatabaseDropRegion{}

// Format implements the NodeFormatter interface.
func (node *AlterDatabaseDropRegion) Format(ctx *FmtCtx) {
	ctx.WriteString("ALTER DATABASE ")
	ctx.FormatNode(&node.Name)
	ctx.WriteString(" DROP REGION ")
	if node.IfExists {
		ctx.WriteString("IF EXISTS ")
	}
	ctx.FormatNode(&node.Region)
}

// AlterDatabasePrimaryRegion represents a ALTER DATABASE PRIMARY REGION ... statement.
type AlterDatabasePrimaryRegion struct {
	Name          Name
	PrimaryRegion Name
}

var _ Statement = &AlterDatabasePrimaryRegion{}

// Format implements the NodeFormatter interface.
func (node *AlterDatabasePrimaryRegion) Format(ctx *FmtCtx) {
	ctx.WriteString("ALTER DATABASE ")
	ctx.FormatNode(&node.Name)
	ctx.WriteString(" PRIMARY REGION ")
	node.PrimaryRegion.Format(ctx)
}

// AlterDatabaseSurvivalGoal represents a ALTER DATABASE SURVIVE ... statement.
type AlterDatabaseSurvivalGoal struct {
	Name         Name
	SurvivalGoal SurvivalGoal
}

var _ Statement = &AlterDatabaseSurvivalGoal{}

// Format implements the NodeFormatter interface.
func (node *AlterDatabaseSurvivalGoal) Format(ctx *FmtCtx) {
	ctx.WriteString("ALTER DATABASE ")
	ctx.FormatNode(&node.Name)
	ctx.WriteString(" ")
	node.SurvivalGoal.Format(ctx)
}

// AlterDatabasePlacement represents a ALTER DATABASE PLACEMENT statement.
type AlterDatabasePlacement struct {
	Name      Name
	Placement DataPlacement
}

var _ Statement = &AlterDatabasePlacement{}

// Format implements the NodeFormatter interface.
func (node *AlterDatabasePlacement) Format(ctx *FmtCtx) {
	ctx.WriteString("ALTER DATABASE ")
	ctx.FormatNode(&node.Name)
	ctx.WriteString(" ")
	node.Placement.Format(ctx)
}

// AlterDatabaseAddSuperRegion represents a
// ALTER DATABASE ADD SUPER REGION ... statement.
type AlterDatabaseAddSuperRegion struct {
	DatabaseName    Name
	SuperRegionName Name
	Regions         []Name
}

var _ Statement = &AlterDatabaseAddSuperRegion{}

// Format implements the NodeFormatter interface.
func (node *AlterDatabaseAddSuperRegion) Format(ctx *FmtCtx) {
	ctx.WriteString("ALTER DATABASE ")
	ctx.FormatNode(&node.DatabaseName)
	ctx.WriteString(" ADD SUPER REGION ")
	ctx.FormatNode(&node.SuperRegionName)
	ctx.WriteString(" VALUES ")
	for i, region := range node.Regions {
		if i != 0 {
			ctx.WriteString(",")
		}
		ctx.FormatNode(&region)
	}
}

// AlterDatabaseDropSuperRegion represents a
// ALTER DATABASE DROP SUPER REGION ... statement.
type AlterDatabaseDropSuperRegion struct {
	DatabaseName    Name
	SuperRegionName Name
}

var _ Statement = &AlterDatabaseDropSuperRegion{}

// Format implements the NodeFormatter interface.
func (node *AlterDatabaseDropSuperRegion) Format(ctx *FmtCtx) {
	ctx.WriteString("ALTER DATABASE ")
	ctx.FormatNode(&node.DatabaseName)
	ctx.WriteString(" DROP SUPER REGION ")
	ctx.FormatNode(&node.SuperRegionName)
}

// AlterDatabaseAlterSuperRegion represents a
// ALTER DATABASE ADD SUPER REGION ... statement.
type AlterDatabaseAlterSuperRegion struct {
	DatabaseName    Name
	SuperRegionName Name
	Regions         []Name
}

var _ Statement = &AlterDatabaseAlterSuperRegion{}

// Format implements the NodeFormatter interface.
func (node *AlterDatabaseAlterSuperRegion) Format(ctx *FmtCtx) {
	ctx.WriteString("ALTER DATABASE ")
	ctx.FormatNode(&node.DatabaseName)
	ctx.WriteString(" ALTER SUPER REGION ")
	ctx.FormatNode(&node.SuperRegionName)
	ctx.WriteString(" VALUES ")
	for i, region := range node.Regions {
		if i != 0 {
			ctx.WriteString(",")
		}
		ctx.FormatNode(&region)
	}
}
