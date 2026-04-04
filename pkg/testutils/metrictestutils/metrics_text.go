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

package metrictestutils

import (
	"bufio"
	"bytes"
	"regexp"
	"sort"
	"strings"

	"github.com/cockroachdb/cockroach/pkg/util/metric"
)

// GetMetricsText scrapes a metrics registry, filters out the metrics according
// to the given regexp, sorts them, and returns them in a multi-line string.
func GetMetricsText(registry *metric.Registry, re *regexp.Regexp) (string, error) {
	ex := metric.MakePrometheusExporter()
	scrape := func(ex *metric.PrometheusExporter) {
		ex.ScrapeRegistry(registry, true /* includeChildMetrics */)
	}
	var in bytes.Buffer
	if err := ex.ScrapeAndPrintAsText(&in, scrape); err != nil {
		return "", err
	}
	sc := bufio.NewScanner(&in)
	var outLines []string
	for sc.Scan() {
		if bytes.HasPrefix(sc.Bytes(), []byte{'#'}) || !re.Match(sc.Bytes()) {
			continue
		}
		outLines = append(outLines, sc.Text())
	}
	if err := sc.Err(); err != nil {
		return "", err
	}
	sort.Strings(outLines)
	metricsText := strings.Join(outLines, "\n")
	return metricsText, nil
}
