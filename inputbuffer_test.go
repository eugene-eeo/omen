package main

import "testing"

func TestInputBuffer(t *testing.T) {
	type testCase struct {
		f func(ib *inputBuffer)
		p int
		s string
	}

	table := []testCase{
		testCase{
			s: "abcd",
			p: 3,
			f: func(ib *inputBuffer) {
				ib.put('a')
				ib.put('b')
				ib.put('d')
				ib.advance(-1)
				ib.put('c')
			},
		},
		testCase{
			s: "",
			p: 0,
			f: func(ib *inputBuffer) {
				ib.put('a')
				ib.put('b')
				ib.advance(-1)
				ib.advance(-1)
				ib.advance(-1)
				ib.delete()
				ib.advance(1)
				ib.backspace()
			},
		},
		testCase{
			s: "aaaa",
			p: 0,
			f: func(ib *inputBuffer) {
				ib.put('a')
				ib.put('a')
				ib.put('a')
				ib.put('a')
				ib.delete()
				ib.advance(-4)
				ib.backspace()
			},
		},
	}

	for i, entry := range table {
		ib := newInputbuffer()
		entry.f(ib)
		if ib.pos != entry.p {
			t.Errorf("[%d], expected ib.pos = %d, got %d instead", i, entry.p, ib.pos)
		}
		if string(ib.buffer) != entry.s {
			t.Errorf("[%d], expected ib.buffer = '%s', got '%s' instead", i, entry.s, string(ib.buffer))
		}
	}
}
