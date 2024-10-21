package pg

import (
	"github.com/jmoiron/sqlx"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

const up = "up"

var migrations = map[string][]string{
	up: {
		//       ↓ - `pg/`.
		//    ↓    - `storage/`.
		// ↓       - `internal/`.
		"../../../db/migrations/00001_tables.up.sql",
	},
}

func migrate(db *sqlx.DB) error {
	_, fp, _, _ := runtime.Caller(0)
	dp := filepath.Dir(fp)
	for _, migration := range migrations[up] {
		f, e := os.Open(filepath.Join(dp, migration))
		if e != nil {
			return e
		}
		b, e := io.ReadAll(f)
		_ = f.Close()
		if e != nil {
			return e
		}
		_, e = db.Exec(string(b))
		if e != nil {
			return e
		}
	}
	return nil
}
