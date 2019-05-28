package main

import "errors"
import "github.com/google/shlex"

var noCommand = errors.New("no command expanded")

type ParsedFormat struct {
	doSub bool
	left  string
	right string
}

func (p *ParsedFormat) Format(query string) string {
	if !p.doSub {
		return p.left
	}
	return p.left + query + p.right
}

func (p *ParsedFormat) Expand(query string) ([]string, error) {
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
	doSub := false
	L := make([]rune, 0, len(cmdFmt))
	R := make([]rune, 0, len(cmdFmt))
	B := &L
	seen_left := false  // prev is a {
	seen_right := false // prev is a }
	for _, r := range cmdFmt {
		switch r {
		case '{':
			// handle '}{'
			if seen_right {
				*B = append(*B, '}')
				seen_right = false
				seen_left = true
				continue
			}
			if seen_left {
				*B = append(*B, '{')
				seen_left = false
				continue
			}
			seen_left = true
			seen_right = false
		case '}':
			if seen_right {
				*B = append(*B, '}')
				seen_right = false
				continue
			}
			// here we've seen '{}', so split here
			if seen_left {
				seen_left = false
				B = &R
				doSub = true
				continue
			}
			seen_left = false
			seen_right = true
		default:
			if seen_left {
				*B = append(*B, '{')
				seen_left = false
			}
			if seen_right {
				*B = append(*B, '}')
				seen_right = false
			}
			*B = append(*B, r)
		}
	}
	// Finally append any excess
	if seen_left {
		*B = append(*B, '{')
	}
	if seen_right {
		*B = append(*B, '}')
	}
	return ParsedFormat{doSub, string(L), string(R)}
}
