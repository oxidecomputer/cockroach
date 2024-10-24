// Copyright 2021 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

import React from "react";
import { Col, Row, Tabs } from "antd";
import { RouteComponentProps } from "react-router-dom";
import classNames from "classnames/bind";
import _ from "lodash";
import { Tooltip } from "antd";
import { Heading } from "@cockroachlabs/ui-components";

import { Anchor } from "src/anchor";
import { Breadcrumbs } from "src/breadcrumbs";
import { CaretRight } from "src/icon/caretRight";
import { StackIcon } from "src/icon/stackIcon";
import { SqlBox } from "src/sql";
import { ColumnDescriptor, SortSetting, SortedTable } from "src/sortedtable";
import {
  SummaryCard,
  SummaryCardItem,
  SummaryCardItemBoolSetting,
} from "src/summaryCard";
import * as format from "src/util/format";
import {
  ascendingAttr,
  columnTitleAttr,
  syncHistory,
  tabAttr,
  tableStatsClusterSetting,
} from "src/util";

import styles from "./databaseTablePage.module.scss";
import { commonStyles } from "src/common";
import { baseHeadingClasses } from "src/transactionsPage/transactionsPageClasses";
import moment, { Moment } from "moment";
import { Search as IndexIcon } from "@cockroachlabs/icons";
import { Link } from "react-router-dom";
import classnames from "classnames/bind";
import booleanSettingStyles from "../settings/booleanSetting.module.scss";
import { DATE_FORMAT_24_UTC } from "src/util/format";
import LoadingError from "../sqlActivity/errorComponent";
import { Loading } from "../loading";
import { UIConfigState } from "../store";
const cx = classNames.bind(styles);
const booleanSettingCx = classnames.bind(booleanSettingStyles);

const { TabPane } = Tabs;

// We break out separate interfaces for some of the nested objects in our data
// so that we can make (typed) test assertions on narrower slices of the data.
//
// The loading and loaded flags help us know when to dispatch the appropriate
// refresh actions.
//
// The overall structure is:
//
//   interface DatabaseTablePageData {
//     databaseName: string;
//     name: string;
//     details: { // DatabaseTablePageDataDetails
//       loading: boolean;
//       loaded: boolean;
//       createStatement: string;
//       replicaCount: number;
//       indexNames: string[];
//       grants: {
//         user: string;
//         privilege: string;
//       }[];
//     };
//     stats: { // DatabaseTablePageDataStats
//       loading: boolean;
//       loaded: boolean;
//       sizeInBytes: number;
//       rangeCount: number;
//       nodesByRegionString: string;
//     };
//     indexStats: { // DatabaseTablePageIndexStats
//       loading: boolean;
//       loaded: boolean;
//       stats: {
//         indexName: string;
//         totalReads: number;
//         lastUsed: Moment;
//         lastUsedType: string;
//       }[];
//       lastReset: Moment;
//     };
//   }
export interface DatabaseTablePageData {
  databaseName: string;
  name: string;
  details: DatabaseTablePageDataDetails;
  stats: DatabaseTablePageDataStats;
  indexStats: DatabaseTablePageIndexStats;
  showNodeRegionsSection?: boolean;
  automaticStatsCollectionEnabled?: boolean;
  hasAdminRole?: UIConfigState["hasAdminRole"];
}

export interface DatabaseTablePageDataDetails {
  loading: boolean;
  loaded: boolean;
  lastError: Error;
  createStatement: string;
  replicaCount: number;
  indexNames: string[];
  grants: Grant[];
  statsLastUpdated?: Moment;
}

export interface DatabaseTablePageIndexStats {
  loading: boolean;
  loaded: boolean;
  lastError: Error;
  stats: IndexStat[];
  lastReset: Moment;
}

interface IndexStat {
  indexName: string;
  totalReads: number;
  lastUsed: Moment;
  lastUsedType: string;
}

interface Grant {
  user: string;
  privilege: string;
}

export interface DatabaseTablePageDataStats {
  loading: boolean;
  loaded: boolean;
  lastError: Error;
  sizeInBytes: number;
  rangeCount: number;
  nodesByRegionString?: string;
}

export interface DatabaseTablePageActions {
  refreshTableDetails: (database: string, table: string) => void;
  refreshTableStats: (database: string, table: string) => void;
  refreshSettings: () => void;
  refreshIndexStats?: (database: string, table: string) => void;
  resetIndexUsageStats?: (database: string, table: string) => void;
  refreshNodes?: () => void;
  refreshUserSQLRoles: () => void;
}

export type DatabaseTablePageProps = DatabaseTablePageData &
  DatabaseTablePageActions &
  RouteComponentProps;

interface DatabaseTablePageState {
  grantSortSetting: SortSetting;
  indexSortSetting: SortSetting;
  tab: string;
}

const indexTabKey = "overview";
const grantsTabKey = "grants";

class DatabaseTableGrantsTable extends SortedTable<Grant> {}

class IndexUsageStatsTable extends SortedTable<IndexStat> {}

export class DatabaseTablePage extends React.Component<
  DatabaseTablePageProps,
  DatabaseTablePageState
> {
  constructor(props: DatabaseTablePageProps) {
    super(props);

    const { history } = this.props;
    const searchParams = new URLSearchParams(history.location.search);
    const currentTab = searchParams.get(tabAttr) || indexTabKey;
    const indexSort: SortSetting = {
      ascending: true,
      columnTitle: "last used",
    };

    const grantSort: SortSetting = {
      ascending: true,
      columnTitle: "username",
    };

    const columnTitle = searchParams.get(columnTitleAttr);
    if (columnTitle) {
      if (currentTab === grantsTabKey) {
        grantSort.columnTitle = columnTitle;
      } else {
        indexSort.columnTitle = columnTitle;
      }
    }

    this.state = {
      indexSortSetting: indexSort,
      grantSortSetting: grantSort,
      tab: currentTab,
    };
  }

  onTabChange = (tab: string): void => {
    this.setState({ ...this.state, tab });

    this.updateUrlAttrFromState(
      tab === grantsTabKey
        ? this.state.grantSortSetting
        : this.state.indexSortSetting,
    );

    syncHistory(
      {
        tab: tab,
      },
      this.props.history,
    );
  };

  componentDidMount(): void {
    this.refresh();
  }

  componentDidUpdate(): void {
    this.refresh();
  }

  private refresh() {
    this.props.refreshUserSQLRoles();
    if (this.props.refreshNodes != null) {
      this.props.refreshNodes();
    }

    if (
      !this.props.details.loaded &&
      !this.props.details.loading &&
      this.props.details.lastError === undefined
    ) {
      return this.props.refreshTableDetails(
        this.props.databaseName,
        this.props.name,
      );
    }

    if (
      !this.props.stats.loaded &&
      !this.props.stats.loading &&
      this.props.stats.lastError === undefined
    ) {
      return this.props.refreshTableStats(
        this.props.databaseName,
        this.props.name,
      );
    }

    if (!this.props.indexStats.loaded && !this.props.indexStats.loading) {
      return this.props.refreshIndexStats(
        this.props.databaseName,
        this.props.name,
      );
    }

    if (this.props.refreshSettings != null) {
      this.props.refreshSettings();
    }
  }

  minDate = moment.utc("0001-01-01"); // minimum value as per UTC

  private changeIndexSortSetting(sortSetting: SortSetting) {
    const stateCopy = { ...this.state };
    stateCopy.indexSortSetting = sortSetting;
    this.setState(stateCopy);
    this.updateUrlAttrFromState(sortSetting);
  }

  private changeGrantSortSetting(sortSetting: SortSetting) {
    const stateCopy = { ...this.state };
    stateCopy.grantSortSetting = sortSetting;
    this.setState(stateCopy);
    this.updateUrlAttrFromState(sortSetting);
  }

  private updateUrlAttrFromState(sortSetting: SortSetting) {
    const { history } = this.props;
    const searchParams = new URLSearchParams(history.location.search);

    searchParams.set(columnTitleAttr, sortSetting.columnTitle);
    searchParams.set(ascendingAttr, String(sortSetting.ascending));
    history.location.search = searchParams.toString();
    history.replace(history.location);
  }

  private getLastResetString() {
    const lastReset = this.props.indexStats.lastReset;
    if (lastReset.isSame(this.minDate)) {
      return "Last reset: Never";
    } else {
      return "Last reset: " + lastReset.format(DATE_FORMAT_24_UTC);
    }
  }

  private getLastUsedString(indexStat: IndexStat) {
    const lastReset = this.props.indexStats.lastReset;
    switch (indexStat.lastUsedType) {
      case "read":
        return indexStat.lastUsed.format("[Last read:] MMM DD, YYYY [at] H:mm");
      case "reset":
      default:
        // TODO(lindseyjin): replace default case with create time after it's added to table_indexes
        if (lastReset.isSame(this.minDate)) {
          return "Never";
        } else {
          return lastReset.format("[Last reset:] MMM DD, YYYY [at] H:mm");
        }
    }
  }

  private indexStatsColumns: ColumnDescriptor<IndexStat>[] = [
    {
      name: "indexes",
      title: "Indexes",
      hideTitleUnderline: true,
      className: cx("index-stats-table__col-indexes"),
      cell: indexStat => (
        <Link
          to={`${this.props.name}/index/${indexStat.indexName}`}
          className={cx("icon__container")}
        >
          <IndexIcon className={cx("icon--s", "icon--primary")} />
          {indexStat.indexName}
        </Link>
      ),
      sort: indexStat => indexStat.indexName,
    },
    {
      name: "total reads",
      title: "Total Reads",
      hideTitleUnderline: true,
      cell: indexStat => format.Count(indexStat.totalReads),
      sort: indexStat => indexStat.totalReads,
    },
    {
      name: "last used",
      title: "Last Used (UTC)",
      hideTitleUnderline: true,
      className: cx("index-stats-table__col-last-used"),
      cell: indexStat => this.getLastUsedString(indexStat),
      sort: indexStat => indexStat.lastUsed,
    },
  ];

  private grantsColumns: ColumnDescriptor<Grant>[] = [
    {
      name: "username",
      title: (
        <Tooltip placement="bottom" title="The user name.">
          User Name
        </Tooltip>
      ),
      cell: grant => grant.user,
      sort: grant => grant.user,
    },
    {
      name: "privilege",
      title: (
        <Tooltip placement="bottom" title="The list of grants for the user.">
          Grants
        </Tooltip>
      ),
      cell: grant => grant.privilege,
      sort: grant => grant.privilege,
    },
  ];

  render(): React.ReactElement {
    const { hasAdminRole } = this.props;
    return (
      <div className="root table-area">
        <section className={baseHeadingClasses.wrapper}>
          <Breadcrumbs
            items={[
              { link: "/databases", name: "Databases" },
              { link: `/database/${this.props.databaseName}`, name: "Tables" },
              {
                link: `/database/${this.props.databaseName}/table/${this.props.name}`,
                name: `Table: ${this.props.name}`,
              },
            ]}
            divider={
              <CaretRight className={cx("icon--xxs", "icon--primary")} />
            }
          />

          <h3
            className={`${baseHeadingClasses.tableName} ${cx(
              "icon__container",
            )}`}
          >
            <StackIcon className={cx("icon--md", "icon--title")} />
            {this.props.name}
          </h3>
        </section>

        <section className={(baseHeadingClasses.wrapper, cx("tab-area"))}>
          <Tabs
            className={commonStyles("cockroach--tabs")}
            onChange={this.onTabChange}
            activeKey={this.state.tab}
          >
            <TabPane
              tab="Overview"
              key={indexTabKey}
              className={cx("tab-pane")}
            >
              <Loading
                loading={this.props.details.loading && this.props.stats.loading}
                page={"table_details"}
                error={
                  this.props.details.lastError || this.props.stats.lastError
                }
                render={() => (
                  <>
                    <Row gutter={18}>
                      <Col className="gutter-row" span={18}>
                        <SqlBox value={this.props.details.createStatement} />
                      </Col>
                    </Row>

                    <Row gutter={18}>
                      <Col className="gutter-row" span={8}>
                        <SummaryCard className={cx("summary-card")}>
                          <SummaryCardItem
                            label="Size"
                            value={format.Bytes(this.props.stats.sizeInBytes)}
                          />
                          <SummaryCardItem
                            label="Replicas"
                            value={this.props.details.replicaCount}
                          />
                          <SummaryCardItem
                            label="Ranges"
                            value={this.props.stats.rangeCount}
                          />
                          {this.props.details.statsLastUpdated && (
                            <SummaryCardItem
                              label="Table Stats Last Updated"
                              value={this.props.details.statsLastUpdated.format(
                                "MMM DD, YYYY [at] H:mm [(UTC)]",
                              )}
                            />
                          )}
                          {this.props.automaticStatsCollectionEnabled !=
                            null && (
                            <SummaryCardItemBoolSetting
                              label="Auto Stats Collection"
                              value={this.props.automaticStatsCollectionEnabled}
                              toolTipText={
                                <span>
                                  {" "}
                                  Automatic statistics can help improve query
                                  performance. Learn how to{" "}
                                  <Anchor
                                    href={tableStatsClusterSetting}
                                    target="_blank"
                                    className={booleanSettingCx(
                                      "crl-hover-text__link-text",
                                    )}
                                  >
                                    manage statistics collection
                                  </Anchor>
                                  .
                                </span>
                              }
                            />
                          )}
                        </SummaryCard>
                      </Col>

                      <Col className="gutter-row" span={10}>
                        <SummaryCard className={cx("summary-card")}>
                          {this.props.showNodeRegionsSection && (
                            <SummaryCardItem
                              label="Regions/Nodes"
                              value={this.props.stats.nodesByRegionString}
                            />
                          )}
                          <SummaryCardItem
                            label="Database"
                            value={this.props.databaseName}
                          />
                          <SummaryCardItem
                            label="Indexes"
                            value={_.join(this.props.details.indexNames, ", ")}
                            className={cx(
                              "database-table-page__indexes--value",
                            )}
                          />
                        </SummaryCard>
                      </Col>
                    </Row>
                    <Row gutter={18}>
                      <SummaryCard
                        className={cx(
                          "summary-card",
                          "index-stats__summary-card",
                        )}
                      >
                        <div className={cx("index-stats__header")}>
                          <Heading type="h5">Index Stats</Heading>
                          <div className={cx("index-stats__reset-info")}>
                            <Tooltip
                              placement="bottom"
                              title="Index stats accumulate from the time the index was created or had its stats reset. Clicking ‘Reset all index stats’ will reset index stats for the entire cluster."
                            >
                              <div
                                className={cx(
                                  "index-stats__last-reset",
                                  "underline",
                                )}
                              >
                                {this.getLastResetString()}
                              </div>
                            </Tooltip>
                            {hasAdminRole && (
                              <div>
                                <a
                                  className={cx(
                                    "action",
                                    "separator",
                                    "index-stats__reset-btn",
                                  )}
                                  onClick={() =>
                                    this.props.resetIndexUsageStats(
                                      this.props.databaseName,
                                      this.props.name,
                                    )
                                  }
                                >
                                  Reset all index stats
                                </a>
                              </div>
                            )}
                          </div>
                        </div>
                        <IndexUsageStatsTable
                          className="index-stats-table"
                          data={this.props.indexStats.stats}
                          columns={this.indexStatsColumns}
                          sortSetting={this.state.indexSortSetting}
                          onChangeSortSetting={this.changeIndexSortSetting.bind(
                            this,
                          )}
                          loading={this.props.indexStats.loading}
                        />
                      </SummaryCard>
                    </Row>
                  </>
                )}
                renderError={() =>
                  LoadingError({
                    statsType: "databases",
                    timeout:
                      this.props.details.lastError?.name
                        ?.toLowerCase()
                        .includes("timeout") ||
                      this.props.stats.lastError?.name
                        ?.toLowerCase()
                        .includes("timeout"),
                  })
                }
              />
            </TabPane>
            <TabPane tab="Grants" key={grantsTabKey} className={cx("tab-pane")}>
              <Loading
                loading={this.props.details.loading}
                page={"table_details_grants"}
                error={this.props.details.lastError}
                render={() => (
                  <DatabaseTableGrantsTable
                    data={this.props.details.grants}
                    columns={this.grantsColumns}
                    sortSetting={this.state.grantSortSetting}
                    onChangeSortSetting={this.changeGrantSortSetting.bind(this)}
                    loading={this.props.details.loading}
                  />
                )}
                renderError={() =>
                  LoadingError({
                    statsType: "databases",
                    timeout: this.props.details.lastError?.name
                      ?.toLowerCase()
                      .includes("timeout"),
                  })
                }
              />
            </TabPane>
          </Tabs>
        </section>
      </div>
    );
  }
}
