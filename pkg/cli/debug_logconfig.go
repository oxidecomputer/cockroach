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

package cli

import (
	"context"
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/util/log/logconfig"
	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var debugCheckLogConfigCmd = &cobra.Command{
	Use:   "check-log-config",
	Short: "test the log config passed via --log",
	RunE:  runDebugCheckLogConfig,
}

var debugLogChanSel logconfig.ChannelList

func runDebugCheckLogConfig(cmd *cobra.Command, args []string) error {
	if err := setupLogging(context.Background(), cmd,
		true /* isServerCmd */, false /* applyconfig */); err != nil {
		return err
	}
	if cliCtx.ambiguousLogDir {
		fmt.Fprintf(stderr, "warning: ambiguous configuration, consider overriding the logging directory\n")
	}

	c := cliCtx.logConfig
	r, err := yaml.Marshal(&c)
	if err != nil {
		return errors.Wrap(err, "printing configuration")
	}

	fmt.Println("# configuration after validation:")
	fmt.Println(string(r))

	_, key := c.Export(debugLogChanSel)
	fmt.Println("# graphical diagram URL:")
	fmt.Printf("http://www.plantuml.com/plantuml/uml/%s\n", key)

	return nil
}
