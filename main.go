package main

import (
	"os"
)

func main() {
	s, e := NewArchiveService()

	if e != nil {
		handleError(&ApplicationError{"setting up archive service", e})
	}

	e = s.ArchiveConversations(os.Args)

	if e != nil {
		handleError(&ApplicationError{"archiving conversations", e})
	}
}
