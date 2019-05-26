package main

import "github.com/gdamore/tcell"
import runewidth "github.com/mattn/go-runewidth"

func unicodeCells(R []rune, width int, fill bool, f func(int, rune, int)) {
	x := 0
	n := len(R)
	for i := 0; x <= width; i++ {
		r := ' '
		if x == width && n > i {
			// if we are at the final width and string is
			// too long then end with ellipsis
			r = '…'
		} else if i < n {
			r = R[i]
		} else if !fill {
			break
		}
		f(x, r, i)
		x += runewidth.RuneWidth(r)
	}
}

func drawPrompt(sc tcell.Screen, ib *inputBuffer, width int) {
	sc.SetContent(0, 0, '>', nil, tcell.StyleDefault.Dim(true))
	m := -1
	unicodeCells(ib.buffer, width-2, true, func(x int, r rune, i int) {
		sc.SetContent(2+x, 0, r, nil, tcell.StyleDefault)
		if i == ib.pos {
			m = x
		}
	})
	if m >= 0 {
		sc.ShowCursor(2+m, 0)
	}
}
