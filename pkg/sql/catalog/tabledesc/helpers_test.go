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

package tabledesc

import (
	"github.com/cockroachdb/cockroach/pkg/clusterversion"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/errors"
)

func ValidatePartitioning(immI catalog.TableDescriptor) error {
	imm, ok := immI.(*immutable)
	if !ok {
		return errors.Errorf("expected immutable descriptor")
	}
	return imm.validatePartitioning()
}

// constraintValidationErrorAccumulator implements catalog.ValidationErrorAccumulator
type constraintValidationErrorAccumulator struct {
	Errors []error
}

// Report implements catalog.ValidationErrorAccumulator
func (cea *constraintValidationErrorAccumulator) Report(err error) {
	cea.Errors = append(cea.Errors, err)
}

// IsActive implements catalog.ValidationErrorAccumulator
func (cea *constraintValidationErrorAccumulator) IsActive(version clusterversion.Key) bool {
	return true
}

func ValidateConstraints(immI catalog.TableDescriptor) error {
	imm, ok := immI.(*immutable)
	if !ok {
		return errors.Errorf("expected immutable descriptor")
	}
	cea := &constraintValidationErrorAccumulator{}
	imm.validateConstraintIDs(cea)
	if cea.Errors == nil {
		return nil
	}
	if len(cea.Errors) > 1 {
		return errors.AssertionFailedf("expected only a single error inside "+
			"validate constraint %q", cea.Errors)
	}
	return cea.Errors[0]
}

func GetPostDeserializationChanges(
	immI catalog.TableDescriptor,
) (catalog.PostDeserializationChanges, error) {
	imm, ok := immI.(*immutable)
	if !ok {
		return catalog.PostDeserializationChanges{}, errors.Errorf("expected immutable descriptor")
	}
	return imm.GetPostDeserializationChanges(), nil
}

var FitColumnToFamily = fitColumnToFamily

func TestingMakeColumn(
	direction descpb.DescriptorMutation_Direction, desc *descpb.ColumnDescriptor,
) catalog.Column {
	return &column{
		maybeMutation: maybeMutation{mutationDirection: direction},
		desc:          desc,
		ordinal:       0,
	}
}
