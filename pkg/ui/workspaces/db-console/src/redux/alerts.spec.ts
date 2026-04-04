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

import { assert } from "chai";
import { Store } from "redux";
import moment from "moment";
import sinon from "sinon";
import { createHashHistory } from "history";

import * as protos from "src/js/protos";
import { cockroach } from "src/js/protos";
import fetchMock from "src/util/fetch-mock";

import { AdminUIState, createAdminUIStore } from "./state";
import {
  AlertLevel,
  alertDataSync,
  staggeredVersionWarningSelector,
  staggeredVersionDismissedSetting,
  disconnectedAlertSelector,
  disconnectedDismissedLocalSetting,
  emailSubscriptionAlertLocalSetting,
  emailSubscriptionAlertSelector,
} from "./alerts";
import { versionsSelector } from "src/redux/nodes";
import { INSTRUCTIONS_BOX_COLLAPSED_KEY, setUIDataKey } from "./uiData";
import {
  livenessReducerObj,
  nodesReducerObj,
  clusterReducerObj,
  healthReducerObj,
} from "./apiReducers";
import MembershipStatus = cockroach.kv.kvserver.liveness.livenesspb.MembershipStatus;

const sandbox = sinon.createSandbox();

describe("alerts", function() {
  let store: Store<AdminUIState>;
  let dispatch: typeof store.dispatch;
  let state: typeof store.getState;

  beforeEach(function() {
    store = createAdminUIStore(createHashHistory());
    dispatch = store.dispatch;
    state = store.getState;
    // localSettings persist values in sessionStorage and
    // this stub disables caching values between tests.
    sandbox.stub(sessionStorage, "getItem").returns(null);
  });

  afterEach(function() {
    sandbox.restore();
    fetchMock.restore();
  });

  describe("selectors", function() {
    describe("versions", function() {
      it("tolerates missing liveness data", function() {
        dispatch(
          nodesReducerObj.receiveData([
            {
              build_info: {
                tag: "0.1",
              },
            },
            {
              build_info: {
                tag: "0.2",
              },
            },
          ]),
        );
        const versions = versionsSelector(state());
        assert.deepEqual(versions, ["0.1", "0.2"]);
      });

      it("ignores decommissioning/decommissioned nodes", function() {
        dispatch(
          nodesReducerObj.receiveData([
            {
              desc: {
                node_id: 1,
              },
              build_info: {
                tag: "0.1",
              },
            },
            {
              desc: {
                node_id: 2,
              },
              build_info: {
                tag: "0.2",
              },
            },
            {
              desc: {
                node_id: 3,
              },
              build_info: {
                tag: "0.3",
              },
            },
          ]),
        );

        dispatch(
          livenessReducerObj.receiveData(
            new protos.cockroach.server.serverpb.LivenessResponse({
              livenesses: [
                {
                  node_id: 1,
                  membership: MembershipStatus.ACTIVE,
                },
                {
                  node_id: 2,
                  membership: MembershipStatus.DECOMMISSIONING,
                },
                {
                  node_id: 3,
                  membership: MembershipStatus.DECOMMISSIONED,
                },
              ],
            }),
          ),
        );

        const versions = versionsSelector(state());
        assert.deepEqual(versions, ["0.1"]);
      });
    });

    describe("version mismatch warning", function() {
      it("requires versions to be loaded before displaying", function() {
        const numAlert = staggeredVersionWarningSelector(state());
        assert.isUndefined(numAlert);
      });

      it("does not display when versions match", function() {
        dispatch(
          nodesReducerObj.receiveData([
            {
              build_info: {
                tag: "0.1",
              },
            },
            {
              build_info: {
                tag: "0.1",
              },
            },
          ]),
        );
        const numAlert = staggeredVersionWarningSelector(state());
        assert.isUndefined(numAlert);
      });

      it("displays when mismatch detected and not dismissed", function() {
        dispatch(
          nodesReducerObj.receiveData([
            {
              // `desc` intentionally omitted (must not affect outcome).
              build_info: {
                tag: "0.1",
              },
            },
            {
              desc: {
                node_id: 1,
              },
              build_info: {
                tag: "0.2",
              },
            },
          ]),
        );
        const numAlert = staggeredVersionWarningSelector(state());
        assert.isObject(numAlert);
        assert.equal(numAlert.level, AlertLevel.WARNING);
        assert.equal(
          numAlert.title,
          "Multiple versions of CockroachDB are running on this cluster.",
        );
      });

      it("does not display if dismissed locally", function() {
        dispatch(
          nodesReducerObj.receiveData([
            {
              build_info: {
                tag: "0.1",
              },
            },
            {
              build_info: {
                tag: "0.2",
              },
            },
          ]),
        );
        dispatch(staggeredVersionDismissedSetting.set(true));
        const numAlert = staggeredVersionWarningSelector(state());
        assert.isUndefined(numAlert);
      });

      it("dismisses by setting local dismissal", function() {
        dispatch(
          nodesReducerObj.receiveData([
            {
              build_info: {
                tag: "0.1",
              },
            },
            {
              build_info: {
                tag: "0.2",
              },
            },
          ]),
        );
        const numAlert = staggeredVersionWarningSelector(state());
        numAlert.dismiss(dispatch, state);
        assert.isTrue(staggeredVersionDismissedSetting.selector(state()));
      });

      it("num alert dismisses by setting local dismissal", function() {
        dispatch(
          nodesReducerObj.receiveData([
            {
              build_info: {
                tag: "0.1",
              },
            },
            {
              build_info: {
                tag: "0.2",
              },
            },
          ]),
        );
        const numAlert = staggeredVersionWarningSelector(state());
        numAlert.dismiss(dispatch, state);
        assert.isTrue(staggeredVersionDismissedSetting.selector(state()));
      });
    });

    describe("disconnected alert", function() {
      it("requires health to be available before displaying", function() {
        const alert = disconnectedAlertSelector(state());
        assert.isUndefined(alert);
      });

      it("does not display when cluster is healthy", function() {
        dispatch(
          healthReducerObj.receiveData(
            new protos.cockroach.server.serverpb.ClusterResponse({}),
          ),
        );
        const alert = disconnectedAlertSelector(state());
        assert.isUndefined(alert);
      });

      it("displays when cluster health endpoint returns an error", function() {
        dispatch(healthReducerObj.errorData(new Error("error")));
        const alert = disconnectedAlertSelector(state());
        assert.isObject(alert);
        assert.equal(alert.level, AlertLevel.CRITICAL);
        assert.equal(
          alert.title,
          "We're currently having some trouble fetching updated data. If this persists, it might be a good idea to check your network connection to the CockroachDB cluster.",
        );
      });

      it("does not display if dismissed locally", function() {
        dispatch(healthReducerObj.errorData(new Error("error")));
        dispatch(disconnectedDismissedLocalSetting.set(moment()));
        const alert = disconnectedAlertSelector(state());
        assert.isUndefined(alert);
      });

      it("dismisses by setting local dismissal", function(done) {
        dispatch(healthReducerObj.errorData(new Error("error")));
        const alert = disconnectedAlertSelector(state());
        const beforeDismiss = moment();

        alert.dismiss(dispatch, state).then(() => {
          assert.isTrue(
            disconnectedDismissedLocalSetting
              .selector(state())
              .isSameOrAfter(beforeDismiss),
          );
          done();
        });
      });
    });

    describe("email signup for release notes alert", () => {
      it("initialized with default 'false' (hidden) state", () => {
        const settingState = emailSubscriptionAlertLocalSetting.selector(
          state(),
        );
        assert.isFalse(settingState);
      });

      it("dismissed by alert#dismiss", async () => {
        // set alert to open state
        dispatch(emailSubscriptionAlertLocalSetting.set(true));
        let openState = emailSubscriptionAlertLocalSetting.selector(state());
        assert.isTrue(openState);

        // dismiss alert
        const alert = emailSubscriptionAlertSelector(state());
        await alert.dismiss(dispatch, state);
        openState = emailSubscriptionAlertLocalSetting.selector(state());
        assert.isFalse(openState);
      });
    });
  });

  describe("data sync listener", function() {
    let sync: () => void;
    beforeEach(function() {
      // We don't care about the responses, we only care that the sync listener
      // is making requests, which can be verified using "inFlight" settings.
      fetchMock.mock({
        matcher: "*",
        method: "GET",
        response: () => 500,
      });

      sync = alertDataSync(store);
    });

    it("dispatches requests for expected data on empty store", function() {
      sync();
      assert.isTrue(state().cachedData.cluster.inFlight);
      assert.isTrue(state().cachedData.nodes.inFlight);
      assert.isTrue(state().cachedData.health.inFlight);
    });

    it("refreshes health function whenever the last health response is no longer valid.", function() {
      dispatch(
        healthReducerObj.receiveData(
          new protos.cockroach.server.serverpb.ClusterResponse({}),
        ),
      );
      dispatch(healthReducerObj.invalidateData());
      sync();
      assert.isTrue(state().cachedData.health.inFlight);
    });

    it("does not do anything when all data is available.", function() {
      dispatch(
        nodesReducerObj.receiveData([
          {
            build_info: {
              tag: "0.1",
            },
          },
        ]),
      );
      dispatch(
        clusterReducerObj.receiveData(
          new protos.cockroach.server.serverpb.ClusterResponse({
            cluster_id: "my-cluster",
          }),
        ),
      );
      dispatch(setUIDataKey(INSTRUCTIONS_BOX_COLLAPSED_KEY, false));
      dispatch(
        healthReducerObj.receiveData(
          new protos.cockroach.server.serverpb.ClusterResponse({}),
        ),
      );

      const expectedState = state();
      sync();
      assert.deepEqual(state(), expectedState);
    });
  });
});
