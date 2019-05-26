package main

import "flag"
import "time"

type cliOptions struct {
	cmdFormat    string
	debounceTime time.Duration
}

func parseFlags() *cliOptions {
	c := &cliOptions{}
	i := int64(50)

	flag.StringVar(&c.cmdFormat, "cmd", "echo '%s'", "command to be executed")
	flag.Int64Var(&i, "debounce", 80, "debounce time (ms)")
	flag.Parse()

	c.debounceTime = time.Duration(i) * time.Millisecond

	return c
}
