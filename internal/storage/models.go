package storage

type LogEntry struct {
	ID        int    `db:"id"`
	Service   string `db:"service"`
	IP        string `db:"ip"`
	Path      string `db:"path"`
	Status    string `db:"status"`
	Method    string `db:"method"`
	IsViewed  bool   `db:"viewed"`
	CreatedAt string `db:"created_at"`
}

type Ban struct {
	ID       int    `db:"id"`
	IP       string `db:"ip"`
	Reason   string `db:"reason"`
	BannedAt string `db:"banned_at"`
}
