package main

import "github.com/gdamore/tcell"
import runewidth "github.com/mattn/go-runewidth"

func unicodeCells(R []rune, width int, fill bool, f func(int, rune)) {
	x := 0
	n := len(R)
	for i := 0; x <= width; i++ {
		r := ' '
		if x == (width-1) && n > i {
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

func drawPrompt(sc tcell.Screen, opt *cliOptions, ib *inputBuffer, width int) {
	xdiff := 0
	unicodeCells([]rune(opt.prompt), width, false, func(x int, r rune) {
		sc.SetContent(x, 0, r, nil, tcell.StyleDefault.Bold(true).Foreground(tcell.ColorRed))
		xdiff = x
	})
	xdiff++
	i := 0
	m := -1
	unicodeCells(ib.buffer, width-xdiff, true, func(x int, r rune) {
		sc.SetContent(xdiff+x, 0, r, nil, tcell.StyleDefault.Bold(true))
		if i == ib.pos {
			m = x
		}
		i++
	})
	if m >= 0 {
		sc.ShowCursor(xdiff+m, 0)
	}
	for i := 0; i < width; i++ {
		sc.SetContent(i, 1, '─', nil, tcell.StyleDefault.Dim(true))
	}
}

func drawPreviewLine(sc tcell.Screen, opt *cliOptions, pd previewLine, width, height int) {
	y := 1 + pd.lineNo
	// don't draw if the y-position is >= height
	if y >= height {
		return
	}
	unicodeCells(pd.line, width, false, func(x int, r rune) {
		sc.SetContent(x, y, r, nil, tcell.StyleDefault)
	})
}
