package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

const (
	conversationsInsertSql = "INSERT INTO conversations(id, title, import_time, create_time, update_time, current_node) VALUES($1, $2, $3, $4, $5, $6) ON CONFLICT (id) DO NOTHING"
	nodesInsertSql         = "INSERT INTO nodes(id, conversation_id, parent, import_time, children) VALUES($1, $2, $3, $4, $5) ON CONFLICT (id) DO NOTHING"
	messagesInsertSql      = "INSERT INTO messages(id, node_id, author_role, author_name, import_time, create_time, update_time, content_type, parts, end_turn, weight, recipient) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) ON CONFLICT (id) DO NOTHING"
)

func GetDatabase() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, e := sql.Open("postgres", psqlInfo)

	if e != nil {
		return nil, e
	}

	return db, nil
}

func CreateConversationTransaction(db *sql.DB) (*sql.Tx, *sql.Stmt, *sql.Stmt, *sql.Stmt, error) {
	t, e := db.Begin()
	if e != nil {
		return nil, nil, nil, nil, e
	}

	cStatement, e := t.Prepare(conversationsInsertSql)
	if e != nil {
		return nil, nil, nil, nil, e
	}

	nStatement, e := t.Prepare(nodesInsertSql)
	if e != nil {
		return nil, nil, nil, nil, e
	}

	mStatement, e := t.Prepare(messagesInsertSql)
	if e != nil {
		return nil, nil, nil, nil, e
	}

	return t, cStatement, nStatement, mStatement, nil
}
