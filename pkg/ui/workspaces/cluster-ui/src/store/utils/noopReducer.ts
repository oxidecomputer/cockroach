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

/**
 * noopReducer is a stub function to use with `createSlice` (@redux-toolkit) as a definition
 * for reducer case which should not change state but it has to define an action which might be
 * handled in sagas for instance.
 *
 * @example
 * ```
 * const slice = createSlice({
 *  name: "someReducer",
 *  reducers: {
 *    someAction: noopReducer,
 *  },
 * });
 *
 * // then it is possible to access this action like this:
 * slice.actions.someAction()
 * ```
 * In this case, action with type "someReducer/someAction" is dispatched, can be handled
 * by middleware but it doesn't change state.
 */
export const noopReducer = (_state: unknown) => {};
