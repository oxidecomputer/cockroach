// Copyright 2018 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package cloud

import (
	"fmt"
	"sort"
	"time"

	"github.com/cockroachdb/cockroach/pkg/roachprod/config"
	"github.com/cockroachdb/cockroach/pkg/roachprod/logger"
	"github.com/cockroachdb/cockroach/pkg/roachprod/vm"
	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
)

type status struct {
	good    []*Cluster
	warn    []*Cluster
	destroy []*Cluster
}

func (s *status) add(c *Cluster, now time.Time) {
	exp := c.ExpiresAt()
	if exp.After(now) {
		if exp.Before(now.Add(2 * time.Hour)) {
			s.warn = append(s.warn, c)
		} else {
			s.good = append(s.good, c)
		}
	} else {
		s.destroy = append(s.destroy, c)
	}
}

// GCClusters checks all cluster to see if they should be deleted. It only
// fails on failure to perform cloud actions. All other actions (load/save
// file, email) do not abort.
func GCClusters(l *logger.Logger, cloud *Cloud, dryrun bool) error {
	now := timeutil.Now()

	var names []string
	for name := range cloud.Clusters {
		if !config.IsLocalClusterName(name) {
			names = append(names, name)
		}
	}
	sort.Strings(names)

	var s status
	users := make(map[string]*status)
	for _, name := range names {
		c := cloud.Clusters[name]
		u := users[c.User]
		if u == nil {
			u = &status{}
			users[c.User] = u
		}
		s.add(c, now)
		u.add(c, now)
	}

	// Compile list of "bad vms" and destroy them.
	var badVMs vm.List
	for _, vm := range cloud.BadInstances {
		// We only delete "bad vms" if they were created more than 1h ago.
		if now.Sub(vm.CreatedAt) >= time.Hour {
			badVMs = append(badVMs, vm)
		}
	}

	if !dryrun {
		if len(badVMs) > 0 {
			// Destroy bad VMs.
			err := vm.FanOut(badVMs, func(p vm.Provider, vms vm.List) error {
				return p.Delete(vms)
			})
			if err != nil {
				fmt.Println(err)
			}
		}

		// Destroy expired clusters.
		for _, c := range s.destroy {
			if err := DestroyCluster(c); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}
