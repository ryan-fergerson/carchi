package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
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
