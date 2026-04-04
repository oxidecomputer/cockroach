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

package nullsink

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/cloud"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/require"
)

func TestNullSinkReadAndWrite(t *testing.T) {
	defer leaktest.AfterTest(t)()

	ctx := context.Background()
	dest := MakeNullSinkStorageURI("foo")

	conf, err := cloud.ExternalStorageConfFromURI(dest, security.RootUserName())
	if err != nil {
		t.Fatal(err)
	}

	s, err := cloud.MakeExternalStorage(ctx, conf, base.ExternalIODirConfig{}, nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	require.Equal(t, roachpb.ExternalStorage{Provider: roachpb.ExternalStorageProvider_null}, s.Conf())
	require.NoError(t, cloud.WriteFile(ctx, s, "", bytes.NewReader([]byte("abc"))))
	sz, err := s.Size(ctx, "")
	require.NoError(t, err)
	require.Equal(t, int64(0), sz)
	_, err = s.ReadFile(ctx, "")
	require.True(t, errors.Is(err, io.EOF))
}
