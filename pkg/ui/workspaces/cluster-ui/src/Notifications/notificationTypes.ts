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

export type NotificationType =
  | "backup-blocked"
  | "command-commit"
  | "expired"
  | "full-table"
  | "network-partition";
export type NotificationSeverity = "low" | "info" | "moderate" | "critical";
export type NotificationTypeProp = {
  key: string;
  title: string;
  description: string;
  severity: NotificationSeverity;
};
export type NotificationProps = {
  id: number;
  read: boolean;
  timestamp: string;
  type: NotificationType;
};

export const notificationTypes: Array<NotificationTypeProp> = [
  {
    key: "backup-blocked",
    title: "Backup blocked on long-running Transaction",
    description:
      "There is a long running transaction that has prevented a backup on a table for more than 1 hour.",
    severity: "low",
  },
  {
    key: "command-commit",
    title: "Command Commit Latency",
    description:
      "Command Commit Latency is > 100ms on at least one node in this cluster. This can result in poor query performance.",
    severity: "low",
  },
  {
    key: "expired",
    title: "Expired License Key",
    description:
      "Your enterprise license key has expired. Enterprise features are disabled until a new license key is set.",
    severity: "moderate",
  },
  {
    key: "full-table",
    title: "Full Table Scan",
    description:
      "There are queries resulting in full table scans in this cluster. Full table scans may result in poor query performance.",
    severity: "low",
  },
  {
    key: "network-partition",
    title: "Network Partition",
    description: "There may be a network partition in this cluster.",
    severity: "critical",
  },
];

export default notificationTypes;
