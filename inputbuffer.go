package main

import "github.com/gdamore/tcell"

type inputBuffer struct {
	pos    int
	buffer []rune
}

func (i *inputBuffer) backspace() {
	// if the buffer is empty, or we are at the start of the buffer
	if len(i.buffer) == 0 || i.pos == 0 {
		return
	}
	i.buffer = append(i.buffer[:i.pos-1], i.buffer[i.pos:]...)
	i.pos--
}

func (i *inputBuffer) advance(direction int) {
	switch direction {
	case -1:
		if i.pos > 0 {
			i.pos--
		}
	case +1:
		if i.pos < len(i.buffer) {
			i.pos++
		}
	}
}

func (i *inputBuffer) put(r rune) {
	if len(i.buffer) <= i.pos {
		i.buffer = append(i.buffer, r)
	} else {
		i.buffer[i.pos] = r
	}
	i.pos++
}

func (i *inputBuffer) handle(ev *tcell.EventKey) (changed bool) {
	switch ev.Key() {
	case tcell.KeyBackspace2:
		fallthrough
	case tcell.KeyBackspace:
		i.backspace()
	case tcell.KeyLeft:
		i.advance(-1)
	case tcell.KeyRight:
		i.advance(+1)
	case tcell.KeyRune:
		i.put(ev.Rune())
	default:
		return false
	}
	return true
}
