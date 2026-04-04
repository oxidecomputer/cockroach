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

package kvserver

import (
	"context"
	"strings"
	"unsafe"

	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/kvserverbase"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/errors"
)

// StoresIterator is the concrete implementation of
// kvserverbase.StoresIterator.
type StoresIterator Stores

var _ kvserverbase.StoresIterator = &StoresIterator{}

// MakeStoresIterator returns a new StoresIterator instance.
func MakeStoresIterator(stores *Stores) *StoresIterator {
	return (*StoresIterator)(stores)
}

// ForEachStore is part of kvserverbase.StoresIterator.
func (s *StoresIterator) ForEachStore(f func(kvserverbase.Store) error) error {
	var err error
	s.storeMap.Range(func(k int64, v unsafe.Pointer) bool {
		store := (*Store)(v)

		err = f((*baseStore)(store))
		return err == nil
	})
	return err
}

// baseStore is the concrete implementation of kvserverbase.Store.
type baseStore Store

var _ kvserverbase.Store = &baseStore{}

// StoreID is part of kvserverbase.Store.
func (s *baseStore) StoreID() roachpb.StoreID {
	store := (*Store)(s)
	return store.StoreID()
}

// Enqueue is part of kvserverbase.Store.
func (s *baseStore) Enqueue(
	ctx context.Context, queue string, rangeID roachpb.RangeID, skipShouldQueue bool,
) error {
	store := (*Store)(s)
	repl, err := store.GetReplica(rangeID)
	if err != nil {
		return err
	}

	_, processErr, enqueueErr := store.Enqueue(ctx, queue, repl, skipShouldQueue, false /* async */)
	if processErr != nil {
		return processErr
	}
	if enqueueErr != nil {
		return enqueueErr
	}
	return nil
}

// SetQueueActive is part of kvserverbase.Store.
func (s *baseStore) SetQueueActive(active bool, queue string) error {
	store := (*Store)(s)
	var kvQueue replicaQueue
	for _, rq := range store.scanner.queues {
		if strings.EqualFold(rq.Name(), queue) {
			kvQueue = rq
			break
		}
	}

	if kvQueue == nil {
		return errors.Errorf("unknown queue %q", queue)
	}

	kvQueue.SetDisabled(!active)
	return nil
}
