package internal

import (
  "fmt"
  "log"
)

type ApplicationError struct {
  Action string
  Err    error
}

func (e *ApplicationError) Error() string {
  return fmt.Sprintf("error during %s: %v", e.Action, e.Err)
}

func HandleError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}
