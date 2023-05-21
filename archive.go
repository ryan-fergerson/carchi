package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type ConversationImporter struct {
	DB           *sql.DB
	DataImporter *DataImporter
	ImportTime   int64
}

func NewArchiveService() (*ConversationImporter, error) {
	db, err := NewDatabaseService().GetDatabaseHandle()

	if err != nil {
		return nil, err
	}

	return &ConversationImporter{
		DB:           db,
		DataImporter: NewDataImporter(),
		ImportTime:   time.Now().Unix(),
	}, nil
}

func (ci *ConversationImporter) ArchiveConversations(args []string) error {
	conversations, err := ci.DataImporter.ImportData(args)

	if err != nil {
		return &ApplicationError{"parsing conversations", err}
	}

	transaction, cStatement, nStatement, mStatement, err := ci.createConversationTransaction(ci.DB)

	if err != nil {
		return &ApplicationError{"starting transaction", err}
	}

	defer func() {
		if err := recover(); err != nil {
			transaction.Rollback()
			fmt.Printf("Recovered in %v", err)
		}
	}()

	err = ci.storeConversations(cStatement, nStatement, mStatement, conversations)

	if err != nil {
		return err
	}

	err = transaction.Commit()

	if err != nil {
		return &ApplicationError{"committing transaction", err}
	}

	return nil
}

func (ci *ConversationImporter) createConversationTransaction(db *sql.DB) (*sql.Tx, *sql.Stmt, *sql.Stmt, *sql.Stmt, error) {
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

func (ci *ConversationImporter) storeConversations(conversationStatement *sql.Stmt, nodeStatement *sql.Stmt, messageStatement *sql.Stmt, conversations []Conversation) error {
	for _, c := range conversations {
		err := ci.storeConversation(conversationStatement, nodeStatement, messageStatement, &c)

		if err != nil {
			return err
		}
	}

	return nil
}

func (ci *ConversationImporter) storeConversation(conversationStatement *sql.Stmt, nodeStatement *sql.Stmt, messageStatement *sql.Stmt, c *Conversation) error {
	if c.Id == "" {
		return &ApplicationError{"creating conversation statement", errors.New("conversation id is required")}
	}

	_, err := conversationStatement.Exec(c.Id, c.Title, ci.ImportTime, c.CreateTime, c.UpdateTime, c.CurrentNode)

	if err != nil {
		return &ApplicationError{"executing conversation statement", err}
	}

	for nodeId, node := range c.Mapping {
		err := ci.storeNode(nodeStatement, messageStatement, c.Id, nodeId, &node)

		if err != nil {
			return err
		}
	}

	return nil
}

func (ci *ConversationImporter) storeNode(nodeStatement *sql.Stmt, messageStatement *sql.Stmt, conversationId, nodeId string, node *Node) error {
	if nodeId == "" {
		return &ApplicationError{"creating node statement", errors.New("node id is required")}
	}

	childrenJson, err := json.Marshal(node.Children)

	if err != nil {
		return &ApplicationError{"marshalling children json", err}
	}

	_, err = nodeStatement.Exec(nodeId, conversationId, node.Parent, ci.ImportTime, childrenJson)

	if err != nil {
		return &ApplicationError{"executing node statement", err}
	}

	if node.Message != nil {
		err = ci.storeMessage(messageStatement, nodeId, node.Message)

		if err != nil {
			return err
		}
	}

	return nil
}

func (ci *ConversationImporter) storeMessage(messageStatement *sql.Stmt, nodeId string, message *NodeMessage) error {
	if message.Id == "" {
		return &ApplicationError{"creating message statement", errors.New("message id is required")}
	}

	partsJson, err := json.Marshal(message.Content.Parts)

	if err != nil {
		return &ApplicationError{"marshalling parts json", err}
	}

	_, err = messageStatement.Exec(
		message.Id,
		nodeId,
		message.Author.Role,
		message.Author.Name,
		ci.ImportTime,
		message.CreateTime,
		message.UpdateTime,
		message.Content.ContentType,
		string(partsJson),
		message.EndTurn,
		message.Weight,
		message.Recipient,
	)

	if err != nil {
		return &ApplicationError{"executing message statement", err}
	}

	return nil
}
