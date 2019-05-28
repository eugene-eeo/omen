package main

import "flag"
import "time"

type cliOptions struct {
	cmdFormat    string
	debounceTime time.Duration
	prompt       string
	allowEmpty   bool
}

func parseFlags() *cliOptions {
	c := &cliOptions{}
	i := int64(0)

	flag.StringVar(&c.cmdFormat, "cmd", "echo '%s'", "command to be executed")
	flag.StringVar(&c.prompt, "prompt", "> ", "prompt string")
	flag.Int64Var(&i, "debounce", 80, "debounce time (ms)")
	flag.BoolVar(&c.allowEmpty, "allowEmpty", false, "allow empty queries")
	flag.Parse()

	c.debounceTime = time.Duration(i) * time.Millisecond

	return c
}
