package main

import (
	"carchi/internal"
	"carchi/web"
	"errors"
	"flag"
)

func main() {
	mode := flag.String("mode", "", "Mode of operation: 'server' or 'archive'")
	flag.Parse()

	if *mode == "" {
		internal.HandleError(&internal.ApplicationError{"main", errors.New("mode argument is required")})
	}

	switch *mode {
	case "archive":
		s, e := internal.NewArchiveService()

		if e != nil {
			internal.HandleError(&internal.ApplicationError{"setting up archive service", e})
		}

		e = s.ArchiveConversations(flag.Args())
	case "server":
		e := web.StartServer()

		if e != nil {
			internal.HandleError(&internal.ApplicationError{"starting web server", e})
		}
	default:
		internal.HandleError(&internal.ApplicationError{"main", errors.New("invalid mode")})
	}
}
