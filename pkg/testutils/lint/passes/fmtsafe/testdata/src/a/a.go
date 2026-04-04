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

package a

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/errors"
	"go.etcd.io/etcd/raft/v3"
)

var unsafeStr = "abc %d"

const constOk = "safe %d"

func init() {
	_ = recover()

	_ = errors.New(unsafeStr) // want `message argument is not a constant expression`

	// Even though the following is trying to opt out of the linter,
	// the opt out fails because the code is not in a test.

	_ = errors.New(unsafeStr /*nolint:fmtsafe*/) // want `message argument is not a constant expression`

	_ = errors.New("safestr")
	_ = errors.New(constOk)
	_ = errors.New("abo" + constOk)
	_ = errors.New("abo" + unsafeStr) // want `message argument is not a constant expression`

	_ = errors.Newf("safe %d", 123)
	_ = errors.Newf(constOk, 123)
	_ = errors.Newf(unsafeStr, 123) // want `format argument is not a constant expression`
	_ = errors.Newf("abo"+constOk, 123)
	_ = errors.Newf("abo"+unsafeStr, 123) // want `format argument is not a constant expression`

	ctx := context.Background()

	log.Errorf(ctx, "safe %d", 123)
	log.Errorf(ctx, constOk, 123)
	log.Errorf(ctx, unsafeStr, 123) // want `format argument is not a constant expression`
	log.Errorf(ctx, "abo"+constOk, 123)
	log.Errorf(ctx, "abo"+unsafeStr, 123) // want `format argument is not a constant expression`

	var m myLogger
	var l raft.Logger = m

	l.Infof("safe %d", 123)
	l.Infof(constOk, 123)
	l.Infof(unsafeStr, 123) // want `format argument is not a constant expression`
	l.Infof("abo"+constOk, 123)
	l.Infof("abo"+unsafeStr, 123) // want `format argument is not a constant expression`
}

type myLogger struct{}

func (m myLogger) Info(args ...interface{}) {
	log.Errorf(context.Background(), "", args...)
}

func (m myLogger) Infof(_ string, args ...interface{}) {
	log.Errorf(context.Background(), "ignoredfmt", args...)
}
