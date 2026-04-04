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

/**
 * Alerts is a collection of selectors which determine if there are any Alerts
 * to display based on the current redux state.
 */

import _ from "lodash";
import moment from "moment";
import { createSelector } from "reselect";
import { Store, Dispatch, Action, AnyAction } from "redux";
import { ThunkAction } from "redux-thunk";

import { LocalSetting } from "./localsettings";
import {
  INSTRUCTIONS_BOX_COLLAPSED_KEY,
  saveUIData,
  loadUIData,
  isInFlight,
  UIDataState,
  UIDataStatus,
} from "./uiData";
import { refreshCluster, refreshNodes, refreshHealth } from "./apiReducers";
import { numNodesByVersionsSelector } from "src/redux/nodes";
import { AdminUIState, AppDispatch } from "./state";

export enum AlertLevel {
  NOTIFICATION,
  WARNING,
  CRITICAL,
  SUCCESS,
}

export interface AlertInfo {
  // Alert Level, which determines visual qualities such as icon and coloring.
  level: AlertLevel;
  // Title to display with the alert.
  title: string;
  // The text of this alert.
  text?: string;
  // Optional hypertext link to be followed when clicking alert.
  link?: string;
}

export interface Alert extends AlertInfo {
  // ThunkAction which will result in this alert being dismissed. This
  // function will be dispatched to the redux store when the alert is dismissed.
  dismiss: ThunkAction<Promise<void>, AdminUIState, void, AnyAction>;
  // Makes alert to be positioned in the top right corner of the screen instead of
  // stretching to full width.
  showAsAlert?: boolean;
  autoClose?: boolean;
  closable?: boolean;
  autoCloseTimeout?: number;
}

const localSettingsSelector = (state: AdminUIState) => state.localSettings;

// Clusterviz Instruction Box collapsed

export const instructionsBoxCollapsedSetting = new LocalSetting(
  INSTRUCTIONS_BOX_COLLAPSED_KEY,
  localSettingsSelector,
  false,
);

const instructionsBoxCollapsedPersistentLoadedSelector = createSelector(
  (state: AdminUIState) => state.uiData,
  (uiData): boolean =>
    uiData &&
    _.has(uiData, INSTRUCTIONS_BOX_COLLAPSED_KEY) &&
    uiData[INSTRUCTIONS_BOX_COLLAPSED_KEY].status === UIDataStatus.VALID,
);

const instructionsBoxCollapsedPersistentSelector = createSelector(
  (state: AdminUIState) => state.uiData,
  (uiData): boolean =>
    uiData &&
    _.has(uiData, INSTRUCTIONS_BOX_COLLAPSED_KEY) &&
    uiData[INSTRUCTIONS_BOX_COLLAPSED_KEY].status === UIDataStatus.VALID &&
    uiData[INSTRUCTIONS_BOX_COLLAPSED_KEY].data,
);

export const instructionsBoxCollapsedSelector = createSelector(
  instructionsBoxCollapsedPersistentLoadedSelector,
  instructionsBoxCollapsedPersistentSelector,
  instructionsBoxCollapsedSetting.selector,
  (persistentLoaded, persistentCollapsed, localSettingCollapsed): boolean => {
    if (persistentLoaded) {
      return persistentCollapsed;
    }
    return localSettingCollapsed;
  },
);

export function setInstructionsBoxCollapsed(collapsed: boolean) {
  return (dispatch: AppDispatch) => {
    dispatch(instructionsBoxCollapsedSetting.set(collapsed));
    dispatch(
      saveUIData({
        key: INSTRUCTIONS_BOX_COLLAPSED_KEY,
        value: collapsed,
      }),
    );
  };
}

////////////////////////////////////////
// Version mismatch.
////////////////////////////////////////
export const staggeredVersionDismissedSetting = new LocalSetting(
  "staggered_version_dismissed",
  localSettingsSelector,
  false,
);

/**
 * Warning when multiple versions of CockroachDB are detected on the cluster.
 * This excludes decommissioned nodes.
 */
export const staggeredVersionWarningSelector = createSelector(
  numNodesByVersionsSelector,
  staggeredVersionDismissedSetting.selector,
  (versionsMap, versionMismatchDismissed): Alert => {
    if (versionMismatchDismissed) {
      return undefined;
    }
    if (!versionsMap || versionsMap.size < 2) {
      return undefined;
    }
    const versionsText = Array.from(versionsMap)
      .map(([k, v]) => `${v} nodes are running on ${k}`)
      .join(" and ")
      .concat(". ");
    return {
      level: AlertLevel.WARNING,
      title: "Multiple versions of CockroachDB are running on this cluster.",
      text:
        versionsText +
        `You can see a list of all nodes and their versions below.
        This may be part of a normal rolling upgrade process, but should be investigated
        if unexpected.`,
      dismiss: (dispatch: AppDispatch) => {
        dispatch(staggeredVersionDismissedSetting.set(true));
        return Promise.resolve();
      },
    };
  },
);

export const disconnectedDismissedLocalSetting = new LocalSetting(
  "disconnected_dismissed",
  localSettingsSelector,
  moment(0),
);

/**
 * Notification when the Admin UI is disconnected from the cluster.
 */
export const disconnectedAlertSelector = createSelector(
  (state: AdminUIState) => state.cachedData.health,
  disconnectedDismissedLocalSetting.selector,
  (health, disconnectedDismissed): Alert => {
    if (!health || !health.lastError) {
      return undefined;
    }

    // Allow local dismissal for one minute.
    const dismissedMaxTime = moment().subtract(1, "m");
    if (disconnectedDismissed.isAfter(dismissedMaxTime)) {
      return undefined;
    }

    return {
      level: AlertLevel.CRITICAL,
      title:
        "We're currently having some trouble fetching updated data. If this persists, it might be a good idea to check your network connection to the CockroachDB cluster.",
      dismiss: (dispatch: Dispatch<Action>) => {
        dispatch(disconnectedDismissedLocalSetting.set(moment()));
        return Promise.resolve();
      },
    };
  },
);

export const emailSubscriptionAlertLocalSetting = new LocalSetting(
  "email_subscription_alert",
  localSettingsSelector,
  false,
);

export const emailSubscriptionAlertSelector = createSelector(
  emailSubscriptionAlertLocalSetting.selector,
  (emailSubscriptionAlert): Alert => {
    if (!emailSubscriptionAlert) {
      return undefined;
    }
    return {
      level: AlertLevel.SUCCESS,
      title: "You successfully signed up for CockroachDB release notes",
      showAsAlert: true,
      autoClose: true,
      closable: false,
      dismiss: (dispatch: Dispatch<Action>) => {
        dispatch(emailSubscriptionAlertLocalSetting.set(false));
        return Promise.resolve();
      },
    };
  },
);

type CreateStatementDiagnosticsAlertPayload = {
  show: boolean;
  status?: "SUCCESS" | "FAILED";
};

export const createStatementDiagnosticsAlertLocalSetting = new LocalSetting<
  AdminUIState,
  CreateStatementDiagnosticsAlertPayload
>("create_stmnt_diagnostics_alert", localSettingsSelector, { show: false });

export const createStatementDiagnosticsAlertSelector = createSelector(
  createStatementDiagnosticsAlertLocalSetting.selector,
  (createStatementDiagnosticsAlert): Alert => {
    if (
      !createStatementDiagnosticsAlert ||
      !createStatementDiagnosticsAlert.show
    ) {
      return undefined;
    }
    const { status } = createStatementDiagnosticsAlert;

    if (status === "FAILED") {
      return {
        level: AlertLevel.CRITICAL,
        title: "There was an error activating statement diagnostics",
        text:
          "Please try activating again. If the problem continues please reach out to customer support.",
        showAsAlert: true,
        dismiss: (dispatch: Dispatch<Action>) => {
          dispatch(
            createStatementDiagnosticsAlertLocalSetting.set({ show: false }),
          );
          return Promise.resolve();
        },
      };
    }
    return {
      level: AlertLevel.SUCCESS,
      title: "Statement diagnostics were successfully activated",
      showAsAlert: true,
      autoClose: true,
      closable: false,
      dismiss: (dispatch: Dispatch<Action>) => {
        dispatch(
          createStatementDiagnosticsAlertLocalSetting.set({ show: false }),
        );
        return Promise.resolve();
      },
    };
  },
);

type CancelStatementDiagnosticsAlertPayload = {
  show: boolean;
  status?: "SUCCESS" | "FAILED";
};

export const cancelStatementDiagnosticsAlertLocalSetting = new LocalSetting<
  AdminUIState,
  CancelStatementDiagnosticsAlertPayload
>("cancel_stmnt_diagnostics_alert", localSettingsSelector, { show: false });

export const cancelStatementDiagnosticsAlertSelector = createSelector(
  cancelStatementDiagnosticsAlertLocalSetting.selector,
  (cancelStatementDiagnosticsAlert): Alert => {
    if (
      !cancelStatementDiagnosticsAlert ||
      !cancelStatementDiagnosticsAlert.show
    ) {
      return undefined;
    }
    const { status } = cancelStatementDiagnosticsAlert;

    if (status === "FAILED") {
      return {
        level: AlertLevel.CRITICAL,
        title: "There was an error cancelling statement diagnostics",
        text:
          "Please try cancelling the statement diagnostic again. If the problem continues please reach out to customer support.",
        showAsAlert: true,
        dismiss: (dispatch: Dispatch<Action>) => {
          dispatch(
            cancelStatementDiagnosticsAlertLocalSetting.set({ show: false }),
          );
          return Promise.resolve();
        },
      };
    }
    return {
      level: AlertLevel.SUCCESS,
      title: "Statement diagnostics were successfully cancelled",
      showAsAlert: true,
      autoClose: true,
      closable: false,
      dismiss: (dispatch: Dispatch<Action>) => {
        dispatch(
          cancelStatementDiagnosticsAlertLocalSetting.set({ show: false }),
        );
        return Promise.resolve();
      },
    };
  },
);

type TerminateSessionAlertPayload = {
  show: boolean;
  status?: "SUCCESS" | "FAILED";
};

export const terminateSessionAlertLocalSetting = new LocalSetting<
  AdminUIState,
  TerminateSessionAlertPayload
>("terminate_session_alert", localSettingsSelector, { show: false });

export const terminateSessionAlertSelector = createSelector(
  terminateSessionAlertLocalSetting.selector,
  (terminateSessionAlert): Alert => {
    if (!terminateSessionAlert || !terminateSessionAlert.show) {
      return undefined;
    }
    const { status } = terminateSessionAlert;

    if (status === "FAILED") {
      return {
        level: AlertLevel.CRITICAL,
        title: "There was an error cancelling the session.",
        text:
          "Please try cancelling again. If the problem continues please reach out to customer support.",
        showAsAlert: true,
        dismiss: (dispatch: Dispatch<Action>) => {
          dispatch(terminateSessionAlertLocalSetting.set({ show: false }));
          return Promise.resolve();
        },
      };
    }
    return {
      level: AlertLevel.SUCCESS,
      title: "Session cancelled.",
      showAsAlert: true,
      autoClose: true,
      closable: false,
      dismiss: (dispatch: Dispatch<Action>) => {
        dispatch(terminateSessionAlertLocalSetting.set({ show: false }));
        return Promise.resolve();
      },
    };
  },
);

type TerminateQueryAlertPayload = {
  show: boolean;
  status?: "SUCCESS" | "FAILED";
};

export const terminateQueryAlertLocalSetting = new LocalSetting<
  AdminUIState,
  TerminateQueryAlertPayload
>("terminate_query_alert", localSettingsSelector, { show: false });

export const terminateQueryAlertSelector = createSelector(
  terminateQueryAlertLocalSetting.selector,
  (terminateQueryAlert): Alert => {
    if (!terminateQueryAlert || !terminateQueryAlert.show) {
      return undefined;
    }
    const { status } = terminateQueryAlert;

    if (status === "FAILED") {
      return {
        level: AlertLevel.CRITICAL,
        title: "There was an error cancelling the statement.",
        text:
          "Please try cancelling again. If the problem continues please reach out to customer support.",
        showAsAlert: true,
        dismiss: (dispatch: Dispatch<Action>) => {
          dispatch(terminateQueryAlertLocalSetting.set({ show: false }));
          return Promise.resolve();
        },
      };
    }
    return {
      level: AlertLevel.SUCCESS,
      title: "Statement cancelled.",
      showAsAlert: true,
      autoClose: true,
      closable: false,
      dismiss: (dispatch: Dispatch<Action>) => {
        dispatch(terminateQueryAlertLocalSetting.set({ show: false }));
        return Promise.resolve();
      },
    };
  },
);

/**
 * Selector which returns an array of all active alerts which should be
 * displayed in the overview list page, these should be non-critical alerts.
 */

export const overviewListAlertsSelector = createSelector(
  staggeredVersionWarningSelector,
  (...alerts: Alert[]): Alert[] => {
    return _.without(alerts, null, undefined);
  },
);

/**
 * Selector which returns an array of all active alerts which should be
 * displayed in the alerts panel, which is embedded within the cluster overview
 * page; currently, this includes all non-critical alerts.
 */
export const panelAlertsSelector = createSelector(
  staggeredVersionWarningSelector,
  (...alerts: Alert[]): Alert[] => {
    return _.without(alerts, null, undefined);
  },
);

/**
 * Selector which returns an array of all active alerts which should be
 * displayed as a banner, which appears at the top of the page and overlaps
 * content in recognition of the severity of the alert; currently, this includes
 * all critical-level alerts.
 */
export const bannerAlertsSelector = createSelector(
  disconnectedAlertSelector,
  emailSubscriptionAlertSelector,
  createStatementDiagnosticsAlertSelector,
  cancelStatementDiagnosticsAlertSelector,
  terminateSessionAlertSelector,
  terminateQueryAlertSelector,
  (...alerts: Alert[]): Alert[] => {
    return _.without(alerts, null, undefined);
  },
);

/**
 * This function, when supplied with a redux store, generates a callback that
 * attempts to populate missing information that has not yet been loaded from
 * the cluster that is needed to show certain alerts. This returned function is
 * intended to be attached to the store as a subscriber.
 */
export function alertDataSync(store: Store<AdminUIState>) {
  const dispatch = store.dispatch as AppDispatch;

  // Memoizers to prevent unnecessary dispatches of alertDataSync if store
  // hasn't changed in an interesting way.
  let lastUIData: UIDataState;

  return () => {
    const state: AdminUIState = store.getState();

    // Always refresh health.
    dispatch(refreshHealth());

    // Load persistent settings which have not yet been loaded.
    const uiData = state.uiData;
    if (uiData !== lastUIData) {
      lastUIData = uiData;
      const keysToMaybeLoad = [INSTRUCTIONS_BOX_COLLAPSED_KEY];
      const keysToLoad = _.filter(keysToMaybeLoad, key => {
        return !(_.has(uiData, key) || isInFlight(state, key));
      });
      if (keysToLoad) {
        dispatch(loadUIData(...keysToLoad));
      }
    }

    // Load Cluster ID once at startup.
    const cluster = state.cachedData.cluster;
    if (cluster && !cluster.data && !cluster.inFlight) {
      dispatch(refreshCluster());
    }

    // Load Nodes initially if it has not yet been loaded.
    const nodes = state.cachedData.nodes;
    if (nodes && !nodes.data && !nodes.inFlight) {
      dispatch(refreshNodes());
    }
  };
}
