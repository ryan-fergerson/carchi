package main

import (
	"archive/zip"
	"encoding/json"
	"io"
	"log"
)

func ReadDataExportFile(args []string) ([]Conversation, error) {
	if len(args) < 2 {
		log.Fatal("You must provide a zip file as an argument.")
	}

	zipFilePath := args[1]

	r, err := zip.OpenReader(zipFilePath)

	if err != nil {
		return nil, err
	}

	return ParseConversations(r)
}

func ParseConversations(r *zip.ReadCloser) ([]Conversation, error) {
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

			if err != nil {
				return nil, err
			}
		}
	}

	return conversations, nil
}
