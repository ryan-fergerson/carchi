package main

import (
	"archive/zip"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func parseConversations(r *zip.ReadCloser) []Conversation {
	var conversations []Conversation

	for _, f := range r.File {
		if f.Name == "conversations.json" {
			rc, err := f.Open()
			handleError(err)
			defer rc.Close()

			content, err := ioutil.ReadAll(rc)
			handleError(err)

			err = json.Unmarshal(content, &conversations)
			handleError(err)

			break
		}
	}

	return conversations
}

func createSqlTransaction(err error, db *sql.DB) (*sql.Tx, *sql.Stmt, *sql.Stmt, *sql.Stmt) {
	tx, err := db.Begin()
	handleError(err)

	stmtConv, err := tx.Prepare(conversationsInsertSql)
	handleError(err)

	stmtNode, err := tx.Prepare(nodesInsertSql)
	handleError(err)

	stmtMsg, err := tx.Prepare(messagesInsertSql)
	handleError(err)

	return tx, stmtConv, stmtNode, stmtMsg
}

func main() {
	// Check if user provided a file path
	if len(os.Args) < 2 {
		log.Fatal("You must provide a zip file as an argument.")
	}

	// Get the file path from the command line argument
	zipFilePath := os.Args[1]

	// Open the zip file
	r, err := zip.OpenReader(zipFilePath)
	handleError(err)
	defer r.Close()

	conversations := parseConversations(r)

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	handleError(err)
	defer db.Close()

	transaction, conversationStatement, nodeStatement, messageStatement := createSqlTransaction(err, db)
	importTime := time.Now().Unix()

	for _, c := range conversations {
		_, err = conversationStatement.Exec(
			c.Id,
			c.Title,
			importTime,
			c.CreateTime,
			c.UpdateTime,
			c.CurrentNode,
		)
		handleError(err)

		for nodeId, node := range c.Mapping {
			childrenJson, err := json.Marshal(node.Children)
			handleError(err)

			_, err = nodeStatement.Exec(
				nodeId,
				c.Id,
				node.Parent,
				importTime,
				childrenJson,
			)
			handleError(err)

			if node.Message != nil {
				partsJson, err := json.Marshal(node.Message.Content.Parts)
				handleError(err)

				_, err = messageStatement.Exec(
					node.Message.Id,
					nodeId,
					node.Message.Author.Role,
					node.Message.Author.Name,
					importTime,
					node.Message.CreateTime,
					node.Message.UpdateTime,
					node.Message.Content.ContentType,
					string(partsJson),
					node.Message.EndTurn,
					node.Message.Weight,
					node.Message.Recipient,
				)
				handleError(err)
			}
		}
	}

	handleError(transaction.Commit())
}
