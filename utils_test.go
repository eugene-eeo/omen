package main

import "reflect"
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
		pf := parseCommandFormat([]rune(entry.cmdFormat))
		cmd, err := pf.Expand(entry.query)
		if err != nil {
			t.Error(n, "Unexpected error", err)
		}
		if len(cmd) != len(entry.expected) {
			t.Error(n, "Expected", entry.expected, "got", cmd)
			continue
		}
		for i, x := range entry.expected {
			if cmd[i] != x {
				t.Error(n, entry.cmdFormat, ":", i, "Expected", x, "got", cmd[i])
			}
		}
	}
}

func TestReplaceCommandFormat(t *testing.T) {
	type tableEntry struct {
		fmt   string
		parts []string
	}
	table := []tableEntry{
		{"ag {}", []string{"ag ", ""}},
		{"ag {{}", []string{"ag {}"}},
		{"ag {{}}", []string{"ag {}"}},
		{"ag {a}", []string{"ag {a}"}},
		{"ag {a}b", []string{"ag {a}b"}},
		{"ag }{", []string{"ag }{"}},
		{"ag }{ {} }{", []string{"ag }{ ", " }{"}},
		{"ag }{ {} }{ {}", []string{"ag }{ ", " }{ ", ""}},
	}
	for _, entry := range table {
		pf := parseCommandFormat([]rune(entry.fmt))
		if !reflect.DeepEqual([]string(pf), entry.parts) {
			t.Error("parseCommandFormat(", entry.fmt, "): expected pf.parts == ", entry.parts, "got", pf)
		}
	}
}
