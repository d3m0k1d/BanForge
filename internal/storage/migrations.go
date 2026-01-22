package storage

const CreateRequestsTable = `
CREATE TABLE IF NOT EXISTS requests (
	id INTEGER PRIMARY KEY,
	service TEXT NOT NULL,
	ip TEXT NOT NULL,
	path TEXT,
	method TEXT,
	status TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_requests_service ON requests(service);
CREATE INDEX IF NOT EXISTS idx_requests_ip ON requests(ip);
CREATE INDEX IF NOT EXISTS idx_requests_status ON requests(status);
CREATE INDEX IF NOT EXISTS idx_requests_created_at ON requests(created_at);
`

// Миграция для bans.db
const CreateBansTable = `
CREATE TABLE IF NOT EXISTS bans (
	id INTEGER PRIMARY KEY,
	ip TEXT UNIQUE NOT NULL,
	reason TEXT,
	banned_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	expired_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_bans_ip ON bans(ip);
`
