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

import React, { useState } from "react";
import { Helmet } from "react-helmet";
import { ArrowLeft } from "@cockroachlabs/icons";
import {
  PlansSortedTable,
  makeExplainPlanColumns,
  PlanHashStats,
} from "./plansTable";
import { Button } from "../../button";
import { SqlBox, SqlBoxSize } from "../../sql";
import { SortSetting } from "../../sortedtable";
import classNames from "classnames/bind";
import styles from "../statementDetails.module.scss";

const cx = classNames.bind(styles);

interface PlanDetailsProps {
  plans: PlanHashStats[];
  sortSetting: SortSetting;
  onChangeSortSetting: (ss: SortSetting) => void;
}

export function PlanDetails({
  plans,
  sortSetting,
  onChangeSortSetting,
}: PlanDetailsProps): React.ReactElement {
  const [plan, setPlan] = useState<PlanHashStats | null>(null);
  const handleDetails = (plan: PlanHashStats): void => {
    setPlan(plan);
  };
  const backToPlanTable = (): void => {
    setPlan(null);
  };

  if (plan) {
    return renderExplainPlan(plan, backToPlanTable);
  } else {
    return renderPlanTable(
      plans,
      handleDetails,
      sortSetting,
      onChangeSortSetting,
    );
  }
}

function renderPlanTable(
  plans: PlanHashStats[],
  handleDetails: (plan: PlanHashStats) => void,
  sortSetting: SortSetting,
  onChangeSortSetting: (ss: SortSetting) => void,
): React.ReactElement {
  const columns = makeExplainPlanColumns(handleDetails);
  return (
    <div className={cx("table-area")}>
      <PlansSortedTable
        columns={columns}
        data={plans}
        className="statements-table"
        sortSetting={sortSetting}
        onChangeSortSetting={onChangeSortSetting}
      />
    </div>
  );
}

function renderExplainPlan(
  plan: PlanHashStats,
  backToPlanTable: () => void,
): React.ReactElement {
  const explainPlan =
    `Plan Gist: ${plan.stats.plan_gists[0]} \n\n` +
    (plan.explain_plan === "" ? "unavailable" : plan.explain_plan);
  return (
    <div>
      <Helmet title="Plan Details" />
      <Button
        onClick={backToPlanTable}
        type="unstyled-link"
        size="small"
        icon={<ArrowLeft fontSize={"10px"} />}
        iconPosition="left"
        className="small-margin"
      >
        All Plans
      </Button>
      <SqlBox value={explainPlan} size={SqlBoxSize.large} />
    </div>
  );
}
