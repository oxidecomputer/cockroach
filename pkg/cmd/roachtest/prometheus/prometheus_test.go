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

package prometheus

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/option"
	logger "github.com/cockroachdb/cockroach/pkg/roachprod/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func nilLogger() *logger.Logger {
	lcfg := logger.Config{
		Stdout: ioutil.Discard,
		Stderr: ioutil.Discard,
	}
	l, err := lcfg.NewLogger("" /* path */)
	if err != nil {
		panic(err)
	}
	return l
}

func TestMakeYAMLConfig(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		desc string

		mockCluster   func(ctrl *gomock.Controller) Cluster
		scrapeConfigs []ScrapeConfig

		expected string
	}{
		{
			desc: "multiple scrape nodes",
			mockCluster: func(ctrl *gomock.Controller) Cluster {
				c := NewMockCluster(ctrl)
				c.EXPECT().
					ExternalIP(ctx, gomock.Any(), []int{1}).
					Return([]string{"127.0.0.1"}, nil)
				c.EXPECT().
					ExternalIP(ctx, gomock.Any(), []int{3, 4, 5}).
					Return([]string{"127.0.0.3", "127.0.0.4", "127.0.0.5"}, nil)
				c.EXPECT().
					ExternalIP(ctx, gomock.Any(), []int{6}).
					Return([]string{"127.0.0.6"}, nil)
				return c
			},
			scrapeConfigs: []ScrapeConfig{
				{
					JobName:     "workload1",
					MetricsPath: "/b",
					ScrapeNodes: []ScrapeNode{
						{
							Nodes: option.NodeListOption([]int{1}),
							Port:  2002,
						},
						{
							Nodes: option.NodeListOption([]int{3, 4, 5}),
							Port:  2003,
						},
					},
				},
				{
					JobName:     "workload2",
					MetricsPath: "/c",
					ScrapeNodes: []ScrapeNode{
						{
							Nodes: option.NodeListOption([]int{6}),
							Port:  2009,
						},
					},
				},
			},
			expected: `global:
  scrape_interval: 10s
  scrape_timeout: 5s
scrape_configs:
- job_name: workload1
  static_configs:
  - targets:
    - 127.0.0.1:2002
    - 127.0.0.3:2003
    - 127.0.0.4:2003
    - 127.0.0.5:2003
  metrics_path: /b
- job_name: workload2
  static_configs:
  - targets:
    - 127.0.0.6:2009
  metrics_path: /c
`,
		},
		{
			desc: "using make commands",
			mockCluster: func(ctrl *gomock.Controller) Cluster {
				c := NewMockCluster(ctrl)
				c.EXPECT().
					ExternalIP(ctx, gomock.Any(), []int{3, 4, 5}).
					Return([]string{"127.0.0.3", "127.0.0.4", "127.0.0.5"}, nil)
				c.EXPECT().
					ExternalIP(ctx, gomock.Any(), []int{6}).
					Return([]string{"127.0.0.6"}, nil)
				c.EXPECT().
					ExternalIP(ctx, gomock.Any(), []int{8, 9}).
					Return([]string{"127.0.0.8", "127.0.0.9"}, nil)
				return c
			},
			scrapeConfigs: []ScrapeConfig{
				MakeWorkloadScrapeConfig(
					"workload",
					[]ScrapeNode{
						{
							Nodes: option.NodeListOption([]int{3, 4, 5}),
							Port:  2005,
						},
						{
							Nodes: option.NodeListOption([]int{6}),
							Port:  2009,
						},
					},
				),
				MakeInsecureCockroachScrapeConfig(
					"roachie",
					option.NodeListOption([]int{8, 9}),
				),
			},
			expected: `global:
  scrape_interval: 10s
  scrape_timeout: 5s
scrape_configs:
- job_name: workload
  static_configs:
  - targets:
    - 127.0.0.3:2005
    - 127.0.0.4:2005
    - 127.0.0.5:2005
    - 127.0.0.6:2009
  metrics_path: /
- job_name: roachie
  static_configs:
  - targets:
    - 127.0.0.8:26258
    - 127.0.0.9:26258
  metrics_path: /_status/vars
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cfg, err := makeYAMLConfig(
				ctx,
				nilLogger(),
				tc.mockCluster(ctrl),
				tc.scrapeConfigs,
			)
			require.NoError(t, err)
			require.Equal(t, tc.expected, cfg)
		})
	}
}
