// Copyright 2020 The Cockroach Authors.
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

import React, {
  forwardRef,
  useCallback,
  useImperativeHandle,
  useState,
} from "react";
import { Modal } from "../modal";
import { Text } from "../text";
import { ICancelQueryRequest } from "src/store/terminateQuery";
export interface TerminateQueryModalRef {
  showModalFor: (req: ICancelQueryRequest) => void;
}

interface TerminateQueryModalProps {
  cancel: (payload: ICancelQueryRequest) => void;
}

const TerminateQueryModal = (
  props: TerminateQueryModalProps,
  ref: React.RefObject<TerminateQueryModalRef>,
) => {
  const { cancel } = props;
  const [visible, setVisible] = useState(false);
  const [req, setReq] = useState<ICancelQueryRequest>();

  const onOkHandler = useCallback(() => {
    cancel(req);
    setVisible(false);
  }, [req, cancel]);

  const onCancelHandler = useCallback(() => setVisible(false), []);

  useImperativeHandle(ref, () => {
    return {
      showModalFor: (r: ICancelQueryRequest) => {
        setReq(r);
        setVisible(true);
      },
    };
  });

  return (
    <Modal
      visible={visible}
      onOk={onOkHandler}
      onCancel={onCancelHandler}
      okText="Yes"
      cancelText="No"
      title="Cancel the Statement"
    >
      <Text>
        Cancelling a statement ends the statement, returning an error to the
        session.
      </Text>
    </Modal>
  );
};

export default forwardRef(TerminateQueryModal);
