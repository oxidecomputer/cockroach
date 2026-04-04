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

package security_test

import (
	"context"
	"testing"
	"time"

	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
)

// TestDummyCreateCACertAndKey is a placeholder for actual testing functions
// TODO(aaron-crl): [tests] write unit tests
func TestDummyCreateCACertAndKey(t *testing.T) {
	defer leaktest.AfterTest(t)()
	_, _, err := security.CreateCACertAndKey(context.Background(), nil, /* loggerFn */
		time.Hour, "test CA cert generation")
	if err != nil {
		t.Fatalf("expected err=nil, got: %s", err)
	}
}

// TestDummyCreateServiceCertAndKey is a placeholder for actual testing functions
// TODO(aaron-crl): [tests] write unit tests
func TestDummyCreateServiceCertAndKey(t *testing.T) {
	defer leaktest.AfterTest(t)()
	caCert, caKey, err := security.CreateCACertAndKey(context.Background(), nil, /* loggerFn */
		time.Hour, "test CA cert generation")
	if err != nil {
		t.Fatalf("expected err=nil, got: %s", err)
	}

	_, _, err = security.CreateServiceCertAndKey(
		context.Background(), nil, /* loggerFn */
		time.Minute,
		"dummy-common-name",
		[]string{"localhost", "127.0.0.1"},
		caCert,
		caKey,
		false, /* serviceCertIsAlsoValidAsClient */
	)
	if err != nil {
		t.Fatalf("expected err=nil, got: %s", err)
	}
}
