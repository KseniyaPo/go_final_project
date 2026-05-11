package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var (
	createShedulerTableSql = `CREATE TABLE scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
    	date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR(256) NOT NULL,
		comment TEXT,
		repeat VARCHAR(128) NOT NULL DEFAULT "")`
	createShedulerDateIndexSql = `CREATE INDEX idx_scheduler_date ON scheduler (date)`
)

func Init(dbFile string) (*sql.DB, error) {
	install := false
	_, err := os.Stat(dbFile)
	if err != nil {
		install = true
	}

	conn, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	if install {
		conn.Exec(createShedulerTableSql)
		conn.Exec(createShedulerDateIndexSql)
	}

	return conn, nil
}
