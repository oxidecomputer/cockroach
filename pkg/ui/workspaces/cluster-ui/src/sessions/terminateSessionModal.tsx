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
import { ICancelSessionRequest } from "src/store/terminateQuery";
import { Modal } from "../modal";
import { Text } from "../text";

export interface TerminateSessionModalRef {
  showModalFor: (req: ICancelSessionRequest) => void;
}

interface TerminateSessionModalProps {
  cancel: (payload: ICancelSessionRequest) => void;
}

const TerminateSessionModal = (
  props: TerminateSessionModalProps,
  ref: React.RefObject<TerminateSessionModalRef>,
) => {
  const { cancel } = props;
  const [visible, setVisible] = useState(false);
  const [req, setReq] = useState<ICancelSessionRequest>();

  const onOkHandler = useCallback(() => {
    cancel(req);
    setVisible(false);
  }, [req, cancel]);

  const onCancelHandler = useCallback(() => setVisible(false), []);

  useImperativeHandle(ref, () => {
    return {
      showModalFor: (r: ICancelSessionRequest) => {
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
      title="Cancel the Session"
    >
      <Text>
        Cancelling a session ends the session, cancelling its associated
        connection. The client that holds this session will receive a
        &quot;connection terminated&quot; event.
      </Text>
    </Modal>
  );
};

export default forwardRef(TerminateSessionModal);
