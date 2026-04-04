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

package rpc

import (
	"crypto/tls"
	"crypto/x509"
	"math"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
)

// GetAddJoinDialOptions returns a standard list of DialOptions for use during
// Add/Join operations.
// TODO(aaron-crl): Possibly fold this into context.go.
func GetAddJoinDialOptions(certPool *x509.CertPool) []grpc.DialOption {
	// Populate the dialOpts.
	var dialOpts []grpc.DialOption

	dialOpts = append(dialOpts, grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(math.MaxInt32),
		grpc.MaxCallSendMsgSize(math.MaxInt32),
	))
	dialOpts = append(dialOpts, grpc.WithDefaultCallOptions(grpc.UseCompressor((snappyCompressor{}).Name())))
	dialOpts = append(dialOpts, grpc.WithNoProxy())
	backoffConfig := backoff.DefaultConfig
	backoffConfig.MaxDelay = maxBackoff
	dialOpts = append(dialOpts, grpc.WithConnectParams(grpc.ConnectParams{
		Backoff:           backoffConfig,
		MinConnectTimeout: minConnectionTimeout}))
	dialOpts = append(dialOpts, grpc.WithKeepaliveParams(clientKeepalive))
	dialOpts = append(dialOpts,
		grpc.WithInitialWindowSize(initialWindowSize),
		grpc.WithInitialConnWindowSize(initialConnWindowSize))

	// Create a tls.Config that allows insecure mode if certPool is not set but
	// requires it if certPool is set.
	var tlsConf tls.Config
	if certPool != nil {
		tlsConf = tls.Config{
			RootCAs: certPool,
		}
	} else {
		// Connect to HTTPS endpoint unverified (effectively HTTP) for CA.
		tlsConf = tls.Config{InsecureSkipVerify: true}
	}

	creds := credentials.NewTLS(&tlsConf)
	dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))

	return dialOpts
}
