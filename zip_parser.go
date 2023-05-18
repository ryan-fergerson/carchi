package main

import (
	"archive/zip"
	"encoding/json"
	"io"
	"log"
)

func ReadZipFile(args []string) (*zip.ReadCloser, error) {
	if len(args) < 2 {
		log.Fatal("You must provide a zip file as an argument.")
	}

	zipFilePath := args[1]

	return zip.OpenReader(zipFilePath)
}

func ParseConversations(r *zip.ReadCloser) []Conversation {
	var conversations []Conversation

	for _, f := range r.File {
		if f.Name == "conversations.json" {
			err := func() error {
				rc, err := f.Open()
				if err != nil {
					return err
				}
				defer rc.Close()

				content, err := io.ReadAll(rc)
				if err != nil {
					return err
				}

				err = json.Unmarshal(content, &conversations)
				if err != nil {
					return err
				}

				return nil
			}()

			handleError(err)
			break
		}
	}

	return conversations
}
