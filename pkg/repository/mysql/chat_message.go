package mysql

import "database/sql"

type ChatMessage struct {
	DB *sql.DB
}
