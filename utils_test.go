package main

import "testing"

func TestExpand(t *testing.T) {
	type tableEntry struct {
		cmdFormat string
		query     string
		expected  []string
	}
	table := []tableEntry{
		{"ag -i -- '{}'", "hello", []string{"ag", "-i", "--", "hello"}},
		{"ag -i -- '{}'", "abc d", []string{"ag", "-i", "--", "abc d"}},
		{"ag -i -- {}", "abc d", []string{"ag", "-i", "--", "abc", "d"}},
	}
	for n, entry := range table {
		cmd, err := expandCommand(entry.cmdFormat, entry.query)
		if err != nil {
			t.Error(n, "Unexpected error", err)
		}
		if len(cmd) != len(entry.expected) {
			t.Error(n, "Expected", entry.expected, "got", cmd)
		}
		for i, x := range entry.expected {
			if x != cmd[i] {
				t.Error(n, entry.cmdFormat, ":", i, "Expected", x, "got", cmd[i])
			}
		}
	}
}

func TestReplaceCommandFormat(t *testing.T) {
	type tableEntry struct {
		fmt   string
		query string
		str   string
	}
	table := []tableEntry{
		{"ag {}", "query", "ag query"},
		{"ag {{}", "query", "ag {}"},
		{"ag {{}}", "query", "ag {}"},
		{"ag {a}", "query", "ag {a}"},
		{"ag {a}b", "query", "ag {a}b"},
		{"ag }{", "query", "ag }{"},
	}
	for _, entry := range table {
		expanded := replaceCommandFormat(entry.fmt, entry.query)
		if expanded != entry.str {
			t.Error("expandCommand(", entry.fmt, ",", entry.query, "): expected", entry.str, "got", expanded)
		}
	}
}
