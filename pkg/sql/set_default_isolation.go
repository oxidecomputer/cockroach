// Copyright 2017 The Cockroach Authors.
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

package sql

import (
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/util/errorutil/unimplemented"
)

func (p *planner) SetSessionCharacteristics(n *tree.SetSessionCharacteristics) (planNode, error) {
	// Note: We also support SET DEFAULT_TRANSACTION_ISOLATION TO ' .... '.
	switch n.Modes.Isolation {
	case tree.SerializableIsolation, tree.UnspecifiedIsolation:
		// Do nothing. All transactions execute with serializable isolation.
	default:
		return nil, pgerror.Newf(pgcode.InvalidParameterValue,
			"unsupported default isolation level: %s", n.Modes.Isolation)
	}

	if err := p.sessionDataMutatorIterator.applyOnEachMutatorError(func(m sessionDataMutator) error {
		// Note: We also support SET DEFAULT_TRANSACTION_PRIORITY TO ' .... '.
		switch n.Modes.UserPriority {
		case tree.UnspecifiedUserPriority:
		default:
			m.SetDefaultTransactionPriority(n.Modes.UserPriority)
		}

		// Note: We also support SET DEFAULT_TRANSACTION_READ_ONLY TO ' .... '.
		switch n.Modes.ReadWriteMode {
		case tree.ReadOnly:
			m.SetDefaultTransactionReadOnly(true)
		case tree.ReadWrite:
			m.SetDefaultTransactionReadOnly(false)
		case tree.UnspecifiedReadWriteMode:
		default:
			return pgerror.Newf(pgcode.InvalidParameterValue,
				"unsupported default read write mode: %s", n.Modes.ReadWriteMode)
		}

		// Note: We also support SET DEFAULT_TRANSACTION_USE_FOLLOWER_READS TO ' .... '.
		//
		// TODO(nvanbenschoten): now that we have a way to set follower_read_timestamp()
		// as the default AS OF SYSTEM TIME value, do we need a way to unset it using
		// the same SET SESSION CHARACTERISTICS AS TRANSACTION mechanism? Currently, the
		// way to do this is SET DEFAULT_TRANSACTION_USE_FOLLOWER_READS TO FALSE;
		if n.Modes.AsOf.Expr != nil {
			if tree.IsFollowerReadTimestampFunction(n.Modes.AsOf, p.semaCtx.SearchPath) {
				m.SetDefaultTransactionUseFollowerReads(true)
			} else {
				return pgerror.Newf(pgcode.InvalidParameterValue,
					"unsupported default as of system time expression, only %s() allowed",
					tree.FollowerReadTimestampFunctionName)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	// Note: We do not support SET DEFAULT_TRANSACTION_DEFERRABLE TO ' ... '.
	switch n.Modes.Deferrable {
	case tree.NotDeferrable, tree.UnspecifiedDeferrableMode:
		// Do nothing. All transactions execute in a NOT DEFERRABLE mode.
	case tree.Deferrable:
		return nil, unimplemented.NewWithIssue(53432, "DEFERRABLE transactions")
	default:
		return nil, pgerror.Newf(pgcode.InvalidParameterValue,
			"unsupported default deferrable mode: %s", n.Modes.Deferrable)
	}

	return newZeroNode(nil /* columns */), nil
}
