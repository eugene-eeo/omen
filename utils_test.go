package main

import "testing"

func TestExpand(t *testing.T) {
	type tableEntry struct {
		cmdFormat string
		query     string
		expected  []string
	}
	table := []tableEntry{
		{"ag -i -- '%s'", "hello", []string{"ag", "-i", "--", "hello"}},
		{"ag -i -- '%s'", "abc d", []string{"ag", "-i", "--", "abc d"}},
		{"ag -i -- %s", "abc d", []string{"ag", "-i", "--", "abc", "d"}},
	}
	for n, entry := range table {
		cmd, err := expandCommand(entry.cmdFormat, entry.query)
		if err != nil {
			t.Error(n, "Unexpected error", err)
		}
		if len(cmd) != len(entry.expected) {
			t.Error(n, "Expected ", entry.expected, " got ", cmd)
		}
		for i, x := range entry.expected {
			if x != cmd[i] {
				t.Error(n, "Expected ", entry.expected, " got ", cmd)
			}
		}
	}
}
