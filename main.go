package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	f, err := ReadZipFile(os.Args)

	if err != nil {
		handleError(&ApplicationError{"reading zip file", err})
	}

	defer f.Close()

	conversations, err := ParseConversations(f)

	if err != nil {
		handleError(&ApplicationError{"parsing conversations", err})
	}

	db, err := GetDatabase()

	if err != nil {
		handleError(&ApplicationError{"getting database", err})
	}

	defer db.Close()

	transaction, conversationStatement, nodeStatement, messageStatement, err := CreateConversationTransaction(db)

	if err != nil {
		handleError(&ApplicationError{"creating conversation transaction", err})
	}

	importTime := time.Now().Unix()

	for _, c := range conversations {
		if c.Id == "" {
			handleError(&ApplicationError{"creating conversation statement",
				errors.New("conversation id is required " + fmt.Sprintf("\n%+v\n", c))})
		}

		_, err = conversationStatement.Exec(
			c.Id,
			c.Title,
			importTime,
			c.CreateTime,
			c.UpdateTime,
			c.CurrentNode,
		)

		if err != nil {
			handleError(&ApplicationError{"executing conversation statement", err})
		}

		for nodeId, node := range c.Mapping {
			if nodeId == "" {
				handleError(&ApplicationError{"creating node statement",
					errors.New("node id is required " + fmt.Sprintf("\n%+v\n", node))})
			}

			childrenJson, err := json.Marshal(node.Children)

			if err != nil {
				handleError(&ApplicationError{"marshalling children json", err})
			}

			_, err = nodeStatement.Exec(
				nodeId,
				c.Id,
				node.Parent,
				importTime,
				childrenJson,
			)

			if err != nil {
				handleError(&ApplicationError{"executing node statement", err})
			}

			if node.Message != nil {
				if node.Message.Id == "" {
					handleError(&ApplicationError{"creating node statement",
						errors.New("message id is required " + fmt.Sprintf("\n%+v\n", node))})
				}

				partsJson, err := json.Marshal(node.Message.Content.Parts)

				if err != nil {
					handleError(&ApplicationError{"marshalling parts json", err})
				}

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

				if err != nil {
					handleError(&ApplicationError{"executing message statement", err})
				}
			}
		}
	}

	err = transaction.Commit()

	if err != nil {
		handleError(&ApplicationError{"committing transaction", err})
	}
}
