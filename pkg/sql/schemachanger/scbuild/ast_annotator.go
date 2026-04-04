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

package scbuild

import (
	"github.com/cockroachdb/cockroach/pkg/sql/parser"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scbuild/internal/scbuildstmt"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/errors"
)

var _ scbuildstmt.TreeAnnotator = (*astAnnotator)(nil)

// astAnnotator creates a copy of the AST that can be annotated or modified
// safely without impacting the original statement.
type astAnnotator struct {
	statement        tree.Statement
	annotation       tree.Annotations
	nonExistentNames map[*tree.TableName]struct{}
}

func newAstAnnotator(original tree.Statement) (*astAnnotator, error) {
	// Clone the original tree by re-parsing the input back into an AST.
	statement, err := parser.ParseOne(original.String())
	if err != nil {
		return nil, err
	}
	return &astAnnotator{
		nonExistentNames: map[*tree.TableName]struct{}{},
		statement:        statement.AST,
		annotation:       tree.MakeAnnotations(statement.NumAnnotations),
	}, nil
}

// GetStatement returns the cloned copy of the AST that is stored inside
// the annotator.
func (ann *astAnnotator) GetStatement() tree.Statement {
	return ann.statement
}

// GetAnnotations implements scbuildstmt.TreeAnnotator.
func (ann *astAnnotator) GetAnnotations() *tree.Annotations {
	// Sanity: Validate the annotations before returning them
	return &ann.annotation
}

// MarkNameAsNonExistent implements scbuildstmt.TreeAnnotator.
func (ann *astAnnotator) MarkNameAsNonExistent(name *tree.TableName) {
	ann.nonExistentNames[name] = struct{}{}
}

// SetUnresolvedNameAnnotation implements scbuildstmt.TreeAnnotator.
func (ann *astAnnotator) SetUnresolvedNameAnnotation(
	unresolvedName *tree.UnresolvedObjectName, annotation interface{},
) {
	unresolvedName.SetAnnotation(&ann.annotation, annotation)
}

// ValidateAnnotations validates if expected modifications have been applied
// on the AST. After the build phase all names should be fully resolved either
// by directly modifying the AST or through annotations.
func (ann *astAnnotator) ValidateAnnotations() {
	// Sanity: Goes through the entire AST and confirms that and table names that
	// appear inside it are fully resolved, meaning that both the database and
	// schema are known.
	f := tree.NewFmtCtx(
		tree.FmtAlwaysQualifyTableNames|tree.FmtMarkRedactionNode,
		tree.FmtAnnotations(&ann.annotation),
		tree.FmtReformatTableNames(func(ctx *tree.FmtCtx, name *tree.TableName) {
			// Name was not found during lookup inside the builder.
			if _, ok := ann.nonExistentNames[name]; ok {
				return
			}
			if name.CatalogName == "" || name.SchemaName == "" {
				panic(errors.AssertionFailedf("unresolved name inside annotated AST "+
					"(%v)", name.String()))
			}
		}))
	f.FormatNode(ann.statement)
}
