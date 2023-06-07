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
		internal.HandleError(&internal.ApplicationError{Action: "checking mode argument", Err: errors.New("mode argument is required")})
	}

	switch *mode {
	case "archive":
		s, e := internal.NewArchiveService()

		if e != nil {
			internal.HandleError(&internal.ApplicationError{Action: "setting up archive service", Err: e})
		}

		e = s.ArchiveConversations(flag.Args())
	case "server":
		e := web.StartServer()

		if e != nil {
			internal.HandleError(&internal.ApplicationError{Action: "starting web server", Err: e})
		}
	default:
		internal.HandleError(&internal.ApplicationError{Action: "starting application", Err: errors.New("invalid mode")})
	}
}
