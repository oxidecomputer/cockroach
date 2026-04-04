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

import moment from "moment";
import * as protos from "@cockroachlabs/crdb-protobuf-client";
import { fromNumber } from "long";

type Timestamp = protos.google.protobuf.ITimestamp;

/**
 * NanoToMilli converts a nanoseconds value into milliseconds.
 */
export function NanoToMilli(nano: number): number {
  return nano / 1.0e6;
}

export function MilliToSeconds(milli: number): number {
  return milli / 1.0e3;
}

/**
 * MilliToNano converts a millisecond value into nanoseconds.
 */
export function MilliToNano(milli: number): number {
  return milli * 1.0e6;
}

/**
 * SecondsToNano converts a second value into nanoseconds.
 */
export function SecondsToNano(sec: number): number {
  return sec * 1.0e9;
}

/**
 * TimestampToMoment converts a Timestamp$Properties object, as seen in wire.proto, to
 * a Moment object. If timestamp is null, it returns the `defaultsIfNull` value which is
 * by default is current time.
 */
export function TimestampToMoment(
  timestamp?: protos.google.protobuf.ITimestamp,
  defaultsIfNull = moment.utc(),
): moment.Moment {
  if (!timestamp) {
    return defaultsIfNull;
  }
  return moment.utc(
    timestamp.seconds.toNumber() * 1e3 + NanoToMilli(timestamp.nanos),
  );
}

/**
 * TimestampToNumber converts a Timestamp$Properties object, as seen in wire.proto, to
 * its unix time. If timestamp is null, it returns the `defaultIfNull` value which is
 * by default is current time.
 */
export function TimestampToNumber(
  timestamp?: protos.google.protobuf.ITimestamp,
  defaultIfNull = moment.utc().unix(),
): number {
  if (!timestamp) {
    return defaultIfNull;
  }
  return timestamp.seconds.toNumber() + NanoToMilli(timestamp.nanos) * 1e-3;
}

/**
 * TimestampToString converts a Timestamp$Properties object, as seen in wire.proto, to
 * its unix time and returns that value as a string. If timestamp is null, it returns
 * the `defaultIfNull` value which is by default is current time.
 */
export function TimestampToString(
  timestamp?: protos.google.protobuf.ITimestamp,
  defaultIfNull = moment.utc().unix(),
): string {
  if (!timestamp) {
    return defaultIfNull.toString();
  }
  return (
    timestamp.seconds.toNumber() +
    NanoToMilli(timestamp.nanos) * 1e-3
  ).toString();
}

/**
 * LongToMoment converts a Long, representing nanos since the epoch, to a Moment
 * object. If timestamp is null, it returns the current time.
 */
export function LongToMoment(timestamp: Long): moment.Moment {
  if (!timestamp) {
    return moment.utc();
  }
  return moment.utc(NanoToMilli(timestamp.toNumber()));
}

/**
 * DurationToNumber converts a Duration object, as seen in wire.proto, to
 * a number representing the duration in seconds. If timestamp is null,
 * it returns the `defaultIfNull` value which is by default 0.
 */
export function DurationToNumber(
  duration?: protos.google.protobuf.IDuration,
  defaultIfNull = 0,
): number {
  if (!duration) {
    return defaultIfNull;
  }
  return duration.seconds.toNumber() + NanoToMilli(duration.nanos) * 1e-3;
}

/**
 * NumberToDuration converts a number representing a duration in seconds
 * to a Duration object.
 */
export function NumberToDuration(
  seconds?: number,
): protos.google.protobuf.IDuration {
  return new protos.google.protobuf.Duration({
    seconds: fromNumber(seconds),
    nanos: SecondsToNano(seconds - Math.floor(seconds)),
  });
}

// durationFromISO8601String function converts a string date in ISO8601 format to moment.Duration
export const durationFromISO8601String = (value: string): moment.Duration => {
  if (!value) {
    return undefined;
  }
  value = value.toUpperCase();
  if (!value.startsWith("P")) {
    value = `PT${value}`;
  }
  return moment.duration(value);
};

export function makeTimestamp(unixTs: number): Timestamp {
  return new protos.google.protobuf.Timestamp({
    seconds: fromNumber(unixTs),
  });
}
