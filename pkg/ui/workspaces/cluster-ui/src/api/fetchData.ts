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

import { cockroach } from "@cockroachlabs/crdb-protobuf-client";
import { RequestError } from "../util";
import { getBasePath } from "./basePath";

interface ProtoBuilder<
  P extends ConstructorType,
  Prop = FirstConstructorParameter<P>,
  R = InstanceType<P>
> {
  new (properties?: Prop): R;
  encode(message: Prop, writer?: protobuf.Writer): protobuf.Writer;
  decode(reader: protobuf.Reader | Uint8Array, length?: number): R;
}

export function toArrayBuffer(encodedRequest: Uint8Array): ArrayBuffer {
  return encodedRequest.buffer.slice(
    encodedRequest.byteOffset,
    encodedRequest.byteOffset + encodedRequest.byteLength,
  );
}

/**
 * @param RespBuilder expects protobuf stub class to build decode response;
 * @param path relative URL path for requested resource;
 * @param ReqBuilder expects protobuf stub to encode request payload. It has to be
 * class type, not instance;
 * @param reqPayload is request payload object;
 * @param timeout is the timeout for the request (optional),
 * format is TimeoutValue (positive integer of at most 8 digits) +
 * TimeoutUnit ( Hour → "H", Minute → "M", Second → "S", Millisecond → "m" ),
 * e.g. "1M" (1 minute), default value "30S" (30 seconds);
 **/
export const fetchData = <P extends ProtoBuilder<P>, T extends ProtoBuilder<T>>(
  RespBuilder: T,
  path: string,
  ReqBuilder?: P,
  reqPayload?: FirstConstructorParameter<P>,
  timeout?: string,
): Promise<InstanceType<T>> => {
  const grpcTimeout = timeout || "30S";
  const params: RequestInit = {
    headers: {
      Accept: "application/x-protobuf",
      "Content-Type": "application/x-protobuf",
      "Grpc-Timeout": grpcTimeout,
    },
    credentials: "same-origin",
  };

  if (reqPayload) {
    const encodedRequest = ReqBuilder.encode(reqPayload).finish();
    params.method = "POST";
    params.body = toArrayBuffer(encodedRequest);
  }
  const basePath = getBasePath();

  return fetch(`${basePath}${path}`, params)
    .then(response => {
      if (!response.ok) {
        return response.arrayBuffer().then(buffer => {
          let respError;
          try {
            respError = cockroach.server.serverpb.ResponseError.decode(
              new Uint8Array(buffer),
            );
          } catch {
            respError = new cockroach.server.serverpb.ResponseError({
              error: response.statusText,
            });
          }
          throw new RequestError(
            response.statusText,
            response.status,
            respError.error,
          );
        });
      }
      return response.arrayBuffer();
    })
    .then(buffer => RespBuilder.decode(new Uint8Array(buffer)));
};

/**
 * fetchDataJSON makes a request for /api/v2 which uses content type JSON.
 * @param path relative path for requested resource.
 * @param reqPayload request payload object.
 */
export function fetchDataJSON<ResponseType, RequestType>(
  path: string,
  reqPayload?: RequestType,
): Promise<ResponseType> {
  const params: RequestInit = {
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
      "X-Cockroach-API-Session": "cookie",
    },
    credentials: "same-origin",
  };

  if (reqPayload) {
    params.method = "POST";
    params.body = JSON.stringify(reqPayload);
  }

  const basePath = getBasePath();

  return fetch(`${basePath}${path}`, params).then(response => {
    if (!response.ok) {
      throw new RequestError(
        response.statusText,
        response.status,
        response.statusText,
      );
    }

    return response.json();
  });
}
