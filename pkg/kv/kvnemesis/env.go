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

package kvnemesis

import (
	"context"
	gosql "database/sql"
	"fmt"
	"time"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/cockroachdb/cockroach/pkg/config/zonepb"
	"github.com/cockroachdb/cockroach/pkg/util/protoutil"
	"github.com/cockroachdb/errors"
)

// Env manipulates the environment (cluster settings, zone configurations) that
// the Applier operates in.
type Env struct {
	sqlDBs []*gosql.DB
}

func (e *Env) anyNode() *gosql.DB {
	// NOTE: There is currently no need to round-robin through the sql gateways,
	// so we always just return the first DB.
	return e.sqlDBs[0]
}

// SetClosedTimestampInterval sets the kv.closed_timestamp.target_duration
// cluster setting to the provided duration.
func (e *Env) SetClosedTimestampInterval(ctx context.Context, d time.Duration) error {
	q := fmt.Sprintf(`SET CLUSTER SETTING kv.closed_timestamp.target_duration = '%s'`, d)
	_, err := e.anyNode().ExecContext(ctx, q)
	return err
}

// ResetClosedTimestampInterval resets the kv.closed_timestamp.target_duration
// cluster setting to its default value.
func (e *Env) ResetClosedTimestampInterval(ctx context.Context) error {
	const q = `SET CLUSTER SETTING kv.closed_timestamp.target_duration TO DEFAULT`
	_, err := e.anyNode().ExecContext(ctx, q)
	return err
}

// UpdateZoneConfig updates the zone configuration with the provided ID using
// the provided function. If no such zone exists, a new one is created.
func (e *Env) UpdateZoneConfig(
	ctx context.Context, zoneID int, updateZone func(*zonepb.ZoneConfig),
) error {
	return crdb.ExecuteTx(ctx, e.anyNode(), nil, func(tx *gosql.Tx) error {
		// Read existing zone configuration.
		var zone zonepb.ZoneConfig
		var zoneRaw []byte
		const q1 = `SELECT config FROM system.zones WHERE id = $1`
		if err := tx.QueryRowContext(ctx, q1, zoneID).Scan(&zoneRaw); err != nil {
			if !errors.Is(err, gosql.ErrNoRows) {
				return err
			}
			// Zone does not exist. Create it.
			zone = zonepb.DefaultZoneConfig()
		} else {
			if err := protoutil.Unmarshal(zoneRaw, &zone); err != nil {
				return errors.Wrap(err, "unmarshaling existing zone")
			}
		}

		// Update zone configuration proto.
		updateZone(&zone)
		zoneRaw, err := protoutil.Marshal(&zone)
		if err != nil {
			return errors.Wrap(err, "marshaling new zone")
		}

		// Rewrite updated zone configuration.
		const q2 = `UPSERT INTO system.zones (id, config) VALUES ($1, $2)`
		_, err = tx.ExecContext(ctx, q2, zoneID, zoneRaw)
		return err
	})
}
