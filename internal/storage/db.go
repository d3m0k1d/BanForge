package storage

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
)

func CreateTables() error {
	// Requests DB
	db_r, err := sql.Open("sqlite",
		"/var/lib/banforge/requests.db?"+
			"mode=rwc&"+
			"_pragma=journal_mode(WAL)&"+
			"_pragma=busy_timeout(30000)&"+
			"_pragma=synchronous(NORMAL)")
	if err != nil {
		return fmt.Errorf("failed to open requests db: %w", err)
	}
	defer db_r.Close()

	_, err = db_r.Exec(CreateRequestsTable)
	if err != nil {
		return fmt.Errorf("failed to create requests table: %w", err)
	}

	// Bans DB
	db_b, err := sql.Open("sqlite",
		"/var/lib/banforge/bans.db?"+
			"mode=rwc&"+
			"_pragma=journal_mode(WAL)&"+
			"_pragma=busy_timeout(30000)&"+
			"_pragma=synchronous(FULL)")
	if err != nil {
		return fmt.Errorf("failed to open bans db: %w", err)
	}
	defer db_b.Close()

	_, err = db_b.Exec(CreateBansTable)
	if err != nil {
		return fmt.Errorf("failed to create bans table: %w", err)
	}
	fmt.Println("Tables created successfully!")
	return nil
}
