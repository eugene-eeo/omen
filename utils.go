package main

import "strings"
import "errors"
import "github.com/google/shlex"

var noCommand = errors.New("no command expanded")

type ParsedFormat []string

func (p ParsedFormat) Format(query string) string {
	return strings.Join(p, query)
}

func (p ParsedFormat) Expand(query string) ([]string, error) {
	s := p.Format(query)
	parts, err := shlex.Split(s)
	if err != nil {
		return nil, err
	}
	if len(parts) == 0 {
		return nil, noCommand
	}
	return parts, nil
}

func parseCommandFormat(cmdFmt []rune) ParsedFormat {
	parts := []string{}
	B := make([]rune, 0, len(cmdFmt))
	seen_left := false  // prev is a {
	seen_right := false // prev is a }
	for _, r := range cmdFmt {
		switch r {
		case '{':
			// handle '}{'
			if seen_right {
				B = append(B, '}')
				seen_right = false
				seen_left = true
				continue
			}
			if seen_left {
				B = append(B, '{')
				seen_left = false
				continue
			}
			seen_left = true
			seen_right = false
		case '}':
			if seen_right {
				B = append(B, '}')
				seen_right = false
				continue
			}
			// here we've seen '{}', so split here
			if seen_left {
				seen_left = false
				parts = append(parts, string(B))
				// reset buffer
				B = B[:0]
				continue
			}
			seen_left = false
			seen_right = true
		default:
			if seen_left {
				B = append(B, '{')
				seen_left = false
			}
			if seen_right {
				B = append(B, '}')
				seen_right = false
			}
			B = append(B, r)
		}
	}
	// Finally append any excess
	if seen_left {
		B = append(B, '{')
	}
	if seen_right {
		B = append(B, '}')
	}
	parts = append(parts, string(B))
	return ParsedFormat(parts)
}
