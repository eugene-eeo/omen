package main

import "fmt"
import "errors"
import "github.com/google/shlex"

var noCommand = errors.New("no command expanded")

func expandCommand(cmdFmt, query string) ([]string, error) {
	s := fmt.Sprintf(cmdFmt, query)
	parts, err := shlex.Split(s)
	if err != nil {
		return nil, err
	}
	if len(parts) == 0 {
		return nil, noCommand
	}
	return parts, nil
}
