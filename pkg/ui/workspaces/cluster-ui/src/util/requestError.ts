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

export class RequestError extends Error {
  status: number;
  constructor(statusText: string, status: number, message?: string) {
    super(statusText);
    this.status = status;
    this.name = "RequestError";
    this.message = message;
  }
}

export function isRequestError(
  error: Error | RequestError,
): error is RequestError {
  return "status" in error && error.name === "RequestError";
}

export function isForbiddenRequestError(error: Error): boolean {
  return isRequestError(error) && error.status === 403; // match to HTTP Forbidden status code
}
