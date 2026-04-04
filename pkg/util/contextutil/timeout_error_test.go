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

package contextutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cockroachdb/errors/errbase"
	"github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
	ctx := context.Background()
	{
		origErr := &TimeoutError{
			operation: "hello",
			timeout:   3 * time.Minute,
			cause:     fmt.Errorf("woo"),
		}
		enc := errbase.EncodeError(ctx, origErr)
		newErr := errbase.DecodeError(ctx, enc)

		assert.Equal(t, origErr.Error(), newErr.Error())
		assert.Equal(t, origErr, newErr)
	}

	{
		origErr := &TimeoutError{
			operation: "hello",
			timeout:   3 * time.Minute,
			took:      4 * time.Minute,
			cause:     fmt.Errorf("woo"),
		}
		enc := errbase.EncodeError(ctx, origErr)
		newErr := errbase.DecodeError(ctx, enc)

		assert.Equal(t, origErr.Error(), newErr.Error())
		assert.Equal(t, origErr, newErr)
	}
}
