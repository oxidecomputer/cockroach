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

package tree

// AlterTenantSetClusterSetting represents an ALTER TENANT
// SET CLUSTER SETTING statement.
type AlterTenantSetClusterSetting struct {
	SetClusterSetting
	TenantID  Expr
	TenantAll bool
}

// Format implements the NodeFormatter interface.
func (n *AlterTenantSetClusterSetting) Format(ctx *FmtCtx) {
	ctx.WriteString("ALTER TENANT ")
	if n.TenantAll {
		ctx.WriteString("ALL")
	} else {
		ctx.FormatNode(n.TenantID)
	}
	ctx.WriteByte(' ')
	ctx.FormatNode(&n.SetClusterSetting)
}

// ShowTenantClusterSetting represents a SHOW CLUSTER SETTING ... FOR TENANT statement.
type ShowTenantClusterSetting struct {
	*ShowClusterSetting
	TenantID Expr
}

// Format implements the NodeFormatter interface.
func (node *ShowTenantClusterSetting) Format(ctx *FmtCtx) {
	ctx.FormatNode(node.ShowClusterSetting)
	ctx.WriteString(" FOR TENANT ")
	ctx.FormatNode(node.TenantID)
}

// ShowTenantClusterSettingList represents a SHOW CLUSTER SETTINGS FOR TENANT statement.
type ShowTenantClusterSettingList struct {
	*ShowClusterSettingList
	TenantID Expr
}

// Format implements the NodeFormatter interface.
func (node *ShowTenantClusterSettingList) Format(ctx *FmtCtx) {
	ctx.FormatNode(node.ShowClusterSettingList)
	ctx.WriteString(" FOR TENANT ")
	ctx.FormatNode(node.TenantID)
}
