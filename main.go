package main

import (
	"encoding/json"
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

func main() {
	f, err := ReadZipFile(os.Args)
	handleError(err)
	defer f.Close()

	conversations := ParseConversations(f)

	db, err := GetDatabase()
	handleError(err)
	defer db.Close()

	transaction, conversationStatement, nodeStatement, messageStatement := CreateConversationTransaction(err, db)
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
