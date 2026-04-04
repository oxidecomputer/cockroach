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

import React, { useCallback, useState } from "react";
import { StatisticType } from "../statsTableUtil/statsTableUtil";
import classNames from "classnames/bind";
import styles from "./sqlActivity.module.scss";
import { Modal } from "../modal";
import { Text } from "../text";

const cx = classNames.bind(styles);

interface clearStatsProps {
  resetSQLStats: () => void;
  tooltipType: StatisticType;
}

const ClearStats = (props: clearStatsProps): React.ReactElement => {
  const [visible, setVisible] = useState(false);
  const onOkHandler = useCallback(() => {
    props.resetSQLStats();
    setVisible(false);
  }, [props]);

  const showModal = (): void => {
    setVisible(true);
  };

  const onCancelHandler = useCallback(() => setVisible(false), []);

  return (
    <>
      <a className={cx("action", "separator")} onClick={showModal}>
        Reset SQL Stats
      </a>
      <Modal
        visible={visible}
        onOk={onOkHandler}
        onCancel={onCancelHandler}
        okText="Continue"
        cancelText="Cancel"
        title="Do you want to reset SQL stats?"
      >
        <Text>
          This action will reset SQL stats on the Statements and Transactions
          pages and crdb_internal tables. Statistics will be cleared and
          unrecoverable for all users across the entire cluster.
        </Text>
      </Modal>
    </>
  );
};

export default ClearStats;
