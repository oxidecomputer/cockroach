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

package certmgr

import (
	"context"
	"crypto/x509"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSelfSignedCert_Err(t *testing.T) {
	ssc := NewSelfSignedCert(-9999, 0, 0, 0)
	require.NotNil(t, ssc)
	require.Nil(t, ssc.Err())
	ssc.Reload(context.Background())
	require.Regexp(t, "cannot represent time as GeneralizedTime", ssc.Err())
	ssc.ClearErr()
	require.Nil(t, ssc.Err())
}

func TestSelfSignedCert_TLSCert(t *testing.T) {
	ssc := NewSelfSignedCert(1, 6, 3, 5*time.Hour)
	require.NotNil(t, ssc)
	require.Nil(t, ssc.Err())
	ssc.Reload(context.Background())
	require.Nil(t, ssc.Err())
	require.NotNil(t, ssc.TLSCert())
	require.Len(t, ssc.TLSCert().Certificate, 1)
	cert, err := x509.ParseCertificate(ssc.TLSCert().Certificate[0])
	require.NoError(t, err)
	expectedUntil := cert.NotBefore.AddDate(1, 6, 3).Add(5 * time.Hour)
	require.Equal(t, expectedUntil, cert.NotAfter)
}
