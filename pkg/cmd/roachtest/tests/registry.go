// Copyright 2018 The Cockroach Authors.
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

package tests

import "github.com/cockroachdb/cockroach/pkg/cmd/roachtest/registry"

// RegisterTests registers all tests to the Registry. This powers `roachtest run`.
func RegisterTests(r registry.Registry) {
	registerAcceptance(r)
	registerActiveRecord(r)
	registerAllocator(r)
	registerAlterPK(r)
	registerAsyncpg(r)
	registerAutoUpgrade(r)
	registerBackup(r)
	registerBackupMixedVersion(r)
	registerBackupNodeShutdown(r)
	registerCancel(r)
	registerClearRange(r)
	registerClockJumpTests(r)
	registerClockMonotonicTests(r)
	registerConnectionLatencyTest(r)
	registerCopy(r)
	registerDecommission(r)
	registerDiskFull(r)
	registerDiskStalledDetection(r)
	registerDjango(r)
	registerDrain(r)
	registerDrop(r)
	registerEncryption(r)
	registerEngineSwitch(r)
	registerFixtures(r)
	registerFlowable(r)
	registerFollowerReads(r)
	registerGopg(r)
	registerGossip(r)
	registerGORM(r)
	registerHibernate(r, hibernateOpts)
	registerHibernate(r, hibernateSpatialOpts)
	registerHotSpotSplits(r)
	registerImportCancellation(r)
	registerImportDecommissioned(r)
	registerImportMixedVersion(r)
	registerImportTPCC(r)
	registerImportTPCH(r)
	registerImportNodeShutdown(r)
	registerInconsistency(r)
	registerIndexes(r)
	registerJasyncSQL(r)
	RegisterJepsen(r)
	registerJobsMixedVersions(r)
	registerKnex(r)
	registerKV(r)
	registerKVContention(r)
	registerKVQuiescenceDead(r)
	registerKVGracefulDraining(r)
	registerKVScalability(r)
	registerKVSplits(r)
	registerKVRangeLookups(r)
	registerKVMultiStoreWithOverload(r)
	registerLargeRange(r)
	registerLedger(r)
	registerLibPQ(r)
	registerLiquibase(r)
	registerNetwork(r)
	registerPebbleWriteThroughput(r)
	registerPebbleYCSB(r)
	registerPgjdbc(r)
	registerPgx(r)
	registerNodeJSPostgres(r)
	registerPop(r)
	registerPsycopg(r)
	registerQueue(r)
	registerQuitAllNodes(r)
	registerQuitTransfersLeases(r)
	registerRebalanceLoad(r)
	registerReplicaGC(r)
	registerRestart(r)
	registerRestoreNodeShutdown(r)
	registerRestore(r)
	registerRoachmart(r)
	registerRoachtest(r)
	registerRubyPG(r)
	registerSchemaChangeBulkIngest(r)
	registerSchemaChangeDatabaseVersionUpgrade(r)
	registerSchemaChangeDuringKV(r)
	registerSchemaChangeIndexTPCC100(r)
	registerSchemaChangeIndexTPCC1000(r)
	registerSchemaChangeDuringTPCC1000(r)
	registerSchemaChangeInvertedIndex(r)
	registerSchemaChangeMixedVersions(r)
	registerSchemaChangeRandomLoad(r)
	registerScrubAllChecksTPCC(r)
	registerScrubIndexOnlyTPCC(r)
	registerSecondaryIndexesMultiVersionCluster(r)
	registerSecure(r)
	registerSequelize(r)
	registerSlowDrain(r)
	registerSQLAlchemy(r)
	registerSQLSmith(r)
	registerSSTableCorruption(r)
	registerSyncTest(r)
	registerSysbench(r)
	registerTLP(r)
	registerTPCC(r)
	registerTPCDSVec(r)
	registerTPCE(r)
	registerTPCHConcurrency(r)
	registerTPCHVec(r)
	registerKVBench(r)
	registerTypeORM(r)
	registerLoadSplits(r)
	registerValidateGrantOption(r)
	registerVersion(r)
	registerYCSB(r)
	registerTPCHBench(r)
	registerOverload(r)
	registerMultiTenantUpgrade(r)
	registerVersionUpgradePublicSchema(r)
	registerRemoveInvalidDatabasePrivileges(r)
}

// RegisterBenchmarks registers all benchmarks to the registry. This powers `roachtest bench`.
//
// TODO(tbg): it's unclear that `roachtest bench` is that useful, perhaps we make everything
// a roachtest but use a `bench` tag to determine what tests to understand as benchmarks.
func RegisterBenchmarks(r registry.Registry) {
	registerIndexesBench(r)
	registerTPCCBench(r)
	registerKVBench(r)
	registerTPCHBench(r)
}
