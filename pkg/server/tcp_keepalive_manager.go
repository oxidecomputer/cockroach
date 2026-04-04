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

package server

import (
	"context"
	"net"
	"sync/atomic"
	"time"

	"github.com/cockroachdb/cmux"
	"github.com/cockroachdb/cockroach/pkg/util/envutil"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

type tcpKeepAliveManager struct {
	// The keepalive duration.
	tcpKeepAlive time.Duration
	// loggedKeepAliveStatus ensures that errors about setting the TCP
	// keepalive status are only reported once.
	loggedKeepAliveStatus int32
}

func makeTCPKeepAliveManager() tcpKeepAliveManager {
	return tcpKeepAliveManager{
		tcpKeepAlive: envutil.EnvOrDefaultDuration("COCKROACH_SQL_TCP_KEEP_ALIVE", time.Minute),
	}
}

// configure attempts to set TCP keep-alive on
// connection. Does not fail on errors.
func (k *tcpKeepAliveManager) configure(ctx context.Context, conn net.Conn) {
	if k.tcpKeepAlive == 0 {
		return
	}

	muxConn, ok := conn.(*cmux.MuxConn)
	if !ok {
		return
	}
	tcpConn, ok := muxConn.Conn.(*net.TCPConn)
	if !ok {
		return
	}

	// Only log success/failure once.
	doLog := atomic.CompareAndSwapInt32(&k.loggedKeepAliveStatus, 0, 1)
	if err := tcpConn.SetKeepAlive(true); err != nil {
		if doLog {
			log.Ops.Warningf(ctx, "failed to enable TCP keep-alive for pgwire: %v", err)
		}
		return

	}
	if err := tcpConn.SetKeepAlivePeriod(k.tcpKeepAlive); err != nil {
		if doLog {
			log.Ops.Warningf(ctx, "failed to set TCP keep-alive duration for pgwire: %v", err)
		}
		return
	}

	if doLog {
		log.VEventf(ctx, 2, "setting TCP keep-alive to %s for pgwire", k.tcpKeepAlive)
	}
}
