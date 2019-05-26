package main

import "github.com/gdamore/tcell"
import runewidth "github.com/mattn/go-runewidth"

// Taken from https://github.com/gdamore/tcell/blob/master/_demos/unicode.go
func puts(s tcell.Screen, style tcell.Style, x, y int, str string) (dx int) {
	i := 0
	var deferred []rune
	dwidth := 0
	zwj := false
	for _, r := range str {
		if r == '\u200d' {
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
			deferred = append(deferred, r)
			zwj = true
			continue
		}
		if zwj {
			deferred = append(deferred, r)
			zwj = false
			continue
		}
		switch runewidth.RuneWidth(r) {
		case 0:
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
		case 1:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 1
		case 2:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 2
		}
		deferred = append(deferred, r)
	}
	if len(deferred) != 0 {
		s.SetContent(x+i, y, deferred[0], deferred[1:], style)
		i += dwidth
	}
	return i
}

func drawPrompt(sc tcell.Screen, ib *inputBuffer) {

	// prompt
	sc.SetContent(0, 0, '>', nil, tcell.StyleDefault.Dim(true))
	sc.SetContent(1, 0, ' ', nil, tcell.StyleDefault)

	// text
	n := len(ib.buffer) + 1
	for i := 0; i < n; i++ {
		style := tcell.StyleDefault
		if i == ib.pos {
			style = style.Reverse(true)
		}
		x := ' '
		if i < n-1 {
			x = ib.buffer[i]
		}
		sc.SetCell(2+i, 0, style, x)
	}

}
