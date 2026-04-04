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

package descmetadata

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/kv"
	"github.com/cockroachdb/cockroach/pkg/settings"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descs"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scexec"
	"github.com/cockroachdb/cockroach/pkg/sql/sessiondata"
	"github.com/cockroachdb/cockroach/pkg/sql/sessiondatapb"
	"github.com/cockroachdb/cockroach/pkg/sql/sessioninit"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlutil"
)

// MetadataUpdaterFactory used to construct a commenter.DescriptorMetadataUpdater, which
// can be used to update comments on schema objects.
type MetadataUpdaterFactory struct {
	ieFactory         sqlutil.SessionBoundInternalExecutorFactory
	collectionFactory *descs.CollectionFactory
	settings          *settings.Values
}

// NewMetadataUpdaterFactory creates a new comment updater factory.
func NewMetadataUpdaterFactory(
	ieFactory sqlutil.SessionBoundInternalExecutorFactory,
	collectionFactory *descs.CollectionFactory,
	settings *settings.Values,
) scexec.DescriptorMetadataUpdaterFactory {
	return MetadataUpdaterFactory{
		ieFactory:         ieFactory,
		collectionFactory: collectionFactory,
		settings:          settings,
	}
}

// NewMetadataUpdater creates a new comment updater, which can be used to
// create / destroy metadata (i.e. comments) associated with different
// schema objects.
func (mf MetadataUpdaterFactory) NewMetadataUpdater(
	ctx context.Context, txn *kv.Txn, sessionData *sessiondata.SessionData,
) scexec.DescriptorMetadataUpdater {
	// Unfortunately, we can't use the session data unmodified, previously the
	// code modifying this metadata would use a circular executor that would ignore
	// any settings set later on. We will intentionally, unset problematic settings
	// here.
	modifiedSessionData := sessionData.Clone()
	modifiedSessionData.ExperimentalDistSQLPlanningMode = sessiondatapb.ExperimentalDistSQLPlanningOn
	return metadataUpdater{
		txn:               txn,
		ie:                mf.ieFactory(ctx, modifiedSessionData),
		collectionFactory: mf.collectionFactory,
		cacheEnabled:      sessioninit.CacheEnabled.Get(mf.settings),
	}
}
