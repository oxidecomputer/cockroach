// Copyright 2019 The Cockroach Authors.
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

// skiperrs connects to a postgres-compatible server with its URL specified as
// the first argument. It then splits stdin into SQL statements and executes
// them on the connection. Errors are printed but do not stop execution.
package main

import (
	gosql "database/sql"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cockroachdb/cockroach/pkg/cmd/cr2pg/sqlstream"
	"github.com/lib/pq"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s <url>\n", os.Args[0])
		os.Exit(1)
	}
	url := os.Args[1]

	connector, err := pq.NewConnector(url)
	if err != nil {
		log.Fatal(err)
	}
	db := gosql.OpenDB(connector)
	defer db.Close()

	stream := sqlstream.NewStream(os.Stdin)
	for {
		stmt, err := stream.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		if _, err := db.Exec(stmt.String()); err != nil {
			fmt.Println(err)
		}
	}
}
