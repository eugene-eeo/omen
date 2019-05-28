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
		doSub bool
		left  string
		right string
	}
	table := []tableEntry{
		{"ag {}", true, "ag ", ""},
		{"ag {{}", false, "ag {}", ""},
		{"ag {{}}", false, "ag {}", ""},
		{"ag {a}", false, "ag {a}", ""},
		{"ag {a}b", false, "ag {a}b", ""},
		{"ag }{", false, "ag }{", ""},
		{"ag }{ {} }{", true, "ag }{ ", " }{"},
	}
	for _, entry := range table {
		pf := parseCommandFormat([]rune(entry.fmt))
		if pf.doSub != entry.doSub {
			t.Error("parseCommandFormat(", entry.fmt, "): expected pf.doSub", entry.doSub, "got", pf.doSub)
		}
		if pf.left != entry.left {
			t.Error("parseCommandFormat(", entry.fmt, "): expected pf.left", entry.left, "got", pf.left)
		}
		if pf.right != entry.right {
			t.Error("parseCommandFormat(", entry.fmt, "): expected pf.right", entry.right, "got", pf.right)
		}
	}
}
