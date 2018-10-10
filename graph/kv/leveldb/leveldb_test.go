// Copyright 2017 The Cayley Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package leveldb

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/caivega/cayley/graph"
	"github.com/caivega/cayley/graph/kv"
	"github.com/caivega/cayley/graph/kv/kvtest"
)

func makeLeveldb(t testing.TB) (kv.BucketKV, graph.Options, func()) {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "cayley_test_"+Type)
	if err != nil {
		t.Fatalf("Could not create working directory: %v", err)
	}
	db, err := Create(tmpDir, nil)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatal("Failed to create Bolt database.", err)
	}
	return db, nil, func() {
		db.Close()
		os.RemoveAll(tmpDir)
	}
}

func TestLeveldb(t *testing.T) {
	kvtest.TestAll(t, makeLeveldb, nil)
}

func BenchmarkLeveldb(b *testing.B) {
	kvtest.BenchmarkAll(b, makeLeveldb, nil)
}
