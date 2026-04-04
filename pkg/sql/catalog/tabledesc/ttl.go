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

package tabledesc

import (
	"time"

	"github.com/cockroachdb/cockroach/pkg/sql/catalog/catpb"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/robfig/cron/v3"
)

// ValidateRowLevelTTL validates that the TTL options are valid.
func ValidateRowLevelTTL(ttl *catpb.RowLevelTTL) error {
	if ttl == nil {
		return nil
	}
	if ttl.DurationExpr == "" {
		return pgerror.Newf(
			pgcode.InvalidParameterValue,
			`"ttl_expire_after" must be set`,
		)
	}
	if ttl.DeleteBatchSize != 0 {
		if err := ValidateTTLBatchSize("ttl_delete_batch_size", ttl.DeleteBatchSize); err != nil {
			return err
		}
	}
	if ttl.SelectBatchSize != 0 {
		if err := ValidateTTLBatchSize("ttl_select_batch_size", ttl.SelectBatchSize); err != nil {
			return err
		}
	}
	if ttl.DeletionCron != "" {
		if err := ValidateTTLCronExpr("ttl_job_cron", ttl.DeletionCron); err != nil {
			return err
		}
	}
	if ttl.RangeConcurrency != 0 {
		if err := ValidateTTLRangeConcurrency("ttl_range_concurrency", ttl.RangeConcurrency); err != nil {
			return err
		}
	}
	if ttl.DeleteRateLimit != 0 {
		if err := ValidateTTLRateLimit("ttl_delete_rate_limit", ttl.DeleteRateLimit); err != nil {
			return err
		}
	}
	if ttl.RowStatsPollInterval != 0 {
		if err := ValidateTTLRowStatsPollInterval("ttl_row_stats_poll_interval", ttl.RowStatsPollInterval); err != nil {
			return err
		}
	}
	return nil
}

// ValidateTTLBatchSize validates the batch size of a TTL.
func ValidateTTLBatchSize(key string, val int64) error {
	if val <= 0 {
		return pgerror.Newf(
			pgcode.InvalidParameterValue,
			`"%s" must be at least 1`,
			key,
		)
	}
	return nil
}

// ValidateTTLRangeConcurrency validates the batch size of a TTL.
func ValidateTTLRangeConcurrency(key string, val int64) error {
	if val <= 0 {
		return pgerror.Newf(
			pgcode.InvalidParameterValue,
			`"%s" must be at least 1`,
			key,
		)
	}
	return nil
}

// ValidateTTLCronExpr validates the cron expression of TTL.
func ValidateTTLCronExpr(key string, str string) error {
	if _, err := cron.ParseStandard(str); err != nil {
		return pgerror.Wrapf(
			err,
			pgcode.InvalidParameterValue,
			`invalid cron expression for "%s"`,
			key,
		)
	}
	return nil
}

// ValidateTTLRowStatsPollInterval validates the automatic statistics field
// of TTL.
func ValidateTTLRowStatsPollInterval(key string, val time.Duration) error {
	if val <= 0 {
		return pgerror.Newf(
			pgcode.InvalidParameterValue,
			`"%s" must be at least 1`,
			key,
		)
	}
	return nil
}

// ValidateTTLRateLimit validates the rate limit parameters of TTL.
func ValidateTTLRateLimit(key string, val int64) error {
	if val <= 0 {
		return pgerror.Newf(
			pgcode.InvalidParameterValue,
			`"%s" must be at least 1`,
			key,
		)
	}
	return nil
}
