// Copyright 2015 The Cockroach Authors.
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
// permissions and limitations under the License. See the AUTHORS file
// for names of contributors.
//
// Author: Marc Berhault (marc@cockroachlabs.com)

package sql_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/cockroachdb/cockroach/client"
	"github.com/cockroachdb/cockroach/security"
	"github.com/cockroachdb/cockroach/security/securitytest"
	"github.com/cockroachdb/cockroach/server"
	"github.com/cockroachdb/cockroach/util/leaktest"
)

func init() {
	security.SetReadFileFn(securitytest.Asset)
}

//go:generate ../util/leaktest/add-leaktest.sh *_test.go

func TestMain(m *testing.M) {
	leaktest.TestMainWithLeakCheck(m)
}

func setup(t *testing.T) (*server.TestServer, *sql.DB, *client.DB) {
	s := server.StartTestServer(nil)
	// SQL requests use "root" which has ALL permissions on everything.
	sqlDB, err := sql.Open("cockroach", fmt.Sprintf("https://%s@%s?certs=test_certs",
		security.RootUser, s.ServingAddr()))
	if err != nil {
		t.Fatal(err)
	}
	// All KV requests need "node" certs.
	kvDB, err := client.Open(fmt.Sprintf("https://%s@%s?certs=test_certs",
		security.NodeUser, s.ServingAddr()))
	if err != nil {
		t.Fatal(err)
	}

	return s, sqlDB, kvDB
}

func cleanup(s *server.TestServer, db *sql.DB) {
	_ = db.Close()
	s.Stop()
}
