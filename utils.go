package main

import "errors"
import "github.com/google/shlex"

var noCommand = errors.New("no command expanded")

func expandCommand(cmdFmt, query string) ([]string, error) {
	s := replaceCommandFormat(cmdFmt, query)
	parts, err := shlex.Split(s)
	if err != nil {
		return nil, err
	}
	if len(parts) == 0 {
		return nil, noCommand
	}
	return parts, nil
}

func replaceCommandFormat(cmdFmt, query string) string {
	R := make([]rune, 0, len(query)+len(cmdFmt))
	seen_left := false  // prev is a {
	seen_right := false // prev is a }
	for _, r := range []rune(cmdFmt) {
		switch r {
		case '{':
			// handle '}{'
			if seen_right {
				R = append(R, '}')
				seen_right = false
				seen_left = true
				continue
			}
			if seen_left {
				R = append(R, '{')
				seen_left = false
				continue
			}
			seen_left = true
			seen_right = false
		case '}':
			if seen_right {
				R = append(R, '}')
				seen_right = false
				continue
			}
			if seen_left {
				// perform substitution
				seen_left = false
				R = append(R, []rune(query)...)
				continue
			}
			seen_left = false
			seen_right = true
		default:
			if seen_left {
				R = append(R, '{')
			}
			if seen_right {
				R = append(R, '}')
			}
			R = append(R, r)
			seen_left = false
			seen_right = false
		}
	}
	// Finally append any excess
	if seen_left {
		R = append(R, '{')
	}
	if seen_right {
		R = append(R, '}')
	}
	return string(R)
}
