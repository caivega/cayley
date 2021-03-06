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

package btree

import (
	"testing"

	"github.com/caivega/cayley/graph"
	"github.com/caivega/cayley/graph/kv"
	"github.com/caivega/cayley/graph/kv/kvtest"
)

func makeBtree(t testing.TB) (kv.BucketKV, graph.Options, func()) {
	return New(), nil, func() {}
}

var conf = &kvtest.Config{
	AlwaysRunIntegration: true,
}

func TestBtree(t *testing.T) {
	kvtest.TestAll(t, makeBtree, conf)
}

func BenchmarkBtree(b *testing.B) {
	kvtest.BenchmarkAll(b, makeBtree, conf)
}
