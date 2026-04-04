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

package execgen

import (
	"github.com/cockroachdb/cockroach/pkg/sql/colexecerror"
	"github.com/cockroachdb/errors"
)

const nonTemplatePanic = "do not call from non-template code"

// Remove unused warnings.
var (
	_ = COPYVAL
	_ = APPENDSLICE
	_ = APPENDVAL
	_ = SETVARIABLESIZE
)

// COPYVAL is a template function that can be used to set a scalar to the value
// of another scalar in such a way that the destination won't be modified if the
// source is.
func COPYVAL(dest, src interface{}) {
	colexecerror.InternalError(errors.AssertionFailedf(nonTemplatePanic))
}

// APPENDSLICE is a template function.
func APPENDSLICE(target, src, destIdx, srcStartIdx, srcEndIdx interface{}) {
	colexecerror.InternalError(errors.AssertionFailedf(nonTemplatePanic))
}

// APPENDVAL is a template function.
func APPENDVAL(target, v interface{}) {
	colexecerror.InternalError(errors.AssertionFailedf(nonTemplatePanic))
}

// SETVARIABLESIZE is a template function.
func SETVARIABLESIZE(target, value interface{}) interface{} {
	colexecerror.InternalError(errors.AssertionFailedf(nonTemplatePanic))
	return nil
}
