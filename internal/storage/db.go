package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

const (
	DBDir     = "/var/lib/banforge/"
	ReqDBPath = DBDir + "requests.db"
	banDBPath = DBDir + "bans.db"
)

var pragmas = map[string]string{
	`journal_mode`: `wal`,
	`synchronous`:  `normal`,
	`busy_timeout`: `30000`,
	// also consider these
	// `temp_store`:   `memory`,
	// `cache_size`:   `1000000000`,
}

func buildSqliteDsn(path string, pragmas map[string]string) string {
	pragmastrs := make([]string, len(pragmas))
	i := 0
	for k, v := range pragmas {
		pragmastrs[i] = (fmt.Sprintf(`pragma=%s(%s)`, k, v))
		i++
	}
	return path + "?" + "mode=rwc&" + strings.Join(pragmastrs, "&")
}

func initDB(dsn, sqlstr string) (err error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return fmt.Errorf("failed to open %q: %w", dsn, err)
	}
	defer func() {
		closeErr := db.Close()
		if closeErr != nil {
			err = errors.Join(err, fmt.Errorf("failed to close %q: %w", dsn, closeErr))
		}
	}()
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return err
}

func CreateTables() (err error) {
	// Requests DB
	err1 := initDB(buildSqliteDsn(ReqDBPath, pragmas), CreateRequestsTable)
	err2 := initDB(buildSqliteDsn(banDBPath, pragmas), CreateBansTable)

	return errors.Join(err1, err2)
}
