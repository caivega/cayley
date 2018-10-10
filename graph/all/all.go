package all

import (
	// supported backends
	_ "github.com/caivega/cayley/graph/kv/bolt"
	_ "github.com/caivega/cayley/graph/kv/btree"
	_ "github.com/caivega/cayley/graph/kv/leveldb"
	_ "github.com/caivega/cayley/graph/memstore"
	_ "github.com/caivega/cayley/graph/nosql/elastic"
	_ "github.com/caivega/cayley/graph/nosql/mongo"
	_ "github.com/caivega/cayley/graph/nosql/ouch"
	_ "github.com/caivega/cayley/graph/sql/cockroach"
	_ "github.com/caivega/cayley/graph/sql/mysql"
	_ "github.com/caivega/cayley/graph/sql/postgres"
)
