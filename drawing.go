package main

import "github.com/gdamore/tcell"
import runewidth "github.com/mattn/go-runewidth"

func unicodeCells(R []rune, width int, fill bool, f func(int, rune)) {
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
		f(x, r)
		x += runewidth.RuneWidth(r)
	}
}

func drawPrompt(sc tcell.Screen, ib *inputBuffer, width int) {
	sc.SetContent(0, 0, '>', nil, tcell.StyleDefault.Dim(true))
	m := ib.pos
	u := ' '
	unicodeCells(ib.buffer, width-2, false, func(x int, r rune) {
		sc.SetContent(2+x, 0, r, nil, tcell.StyleDefault)
		if x == m {
			u = r
		}
	})
	sc.SetContent(2+m, 0, u, nil, tcell.StyleDefault.Reverse(true))
}
