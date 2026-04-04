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
  useRef,
  useEffect,
  forwardRef,
  useImperativeHandle,
} from "react";

type FileTypes = "text/plain" | "application/json";

export interface DownloadAsFileProps {
  fileName?: string;
  fileType?: FileTypes;
  content?: string;
}

export interface DownloadFileRef {
  download: (name: string, type: FileTypes, body: string) => void;
}

/*
 * DownloadFile can download file in two modes `default` and `imperative`.
 * `Default` mode - when DownloadFile wraps component which should trigger
 * downloading and can work only if content of file is already available.
 *
 * For example:
 * ```
 * <DownloadFile fileName="example.txt" fileType="text/plain" content="Some text">
 *   <button>Download</download>
 * </DownloadFile>
 * ```
 *
 * `Imperative` mode allows initiate file download in async way, and trigger
 * download manually.
 *
 * For example:
 * ```
 * downloadRef = React.createRef<DownloadFileRef>();
 *
 * fetchData = () => {
 *   Promise.resolve().then((someText) =>
 *     this.downloadRef.current.download("example.txt", "text/plain", someText))
 * }
 *
 * <DownloadFile ref={downloadRef} />
 * <button onClick={fetchData}>Download</button>
 * ```
 * */
// tslint:disable-next-line:variable-name
export const DownloadFile = forwardRef<DownloadFileRef, DownloadAsFileProps>(
  (props, ref) => {
    const { children, fileName, fileType, content } = props;
    const anchorRef = useRef<HTMLAnchorElement>();

    const bootstrapFile = (name: string, type: FileTypes, body: string) => {
      const anchorElement = anchorRef.current;
      const file = new Blob([body], { type });
      anchorElement.href = URL.createObjectURL(file);
      anchorElement.download = name;
    };

    useEffect(() => {
      if (content === undefined) {
        return;
      }
      bootstrapFile(fileName, fileType, content);
    }, [fileName, fileType, content]);

    useImperativeHandle(ref, () => ({
      download: (name: string, type: FileTypes, body: string) => {
        bootstrapFile(name, type, body);
        anchorRef.current.click();
      },
    }));

    return <a ref={anchorRef}>{children}</a>;
  },
);
