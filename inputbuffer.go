package main

import "github.com/gdamore/tcell"

type inputBuffer struct {
	pos    int
	buffer []rune
}

func newInputbuffer() *inputBuffer {
	return &inputBuffer{
		pos:    0,
		buffer: make([]rune, 0, 50),
	}
}

func (i *inputBuffer) delete() bool {
	// if we're at the end of the buffer
	if i.pos == len(i.buffer) {
		return false
	}
	i.buffer = append(i.buffer[:i.pos], i.buffer[i.pos+1:]...)
	return true
}

func (i *inputBuffer) backspace() bool {
	// if the buffer is empty, or we are at the start of the buffer
	if len(i.buffer) == 0 || i.pos == 0 {
		return false
	}
	i.buffer = append(i.buffer[:i.pos-1], i.buffer[i.pos:]...)
	i.pos--
	return true
}

func (i *inputBuffer) advance(direction int) {
	max := len(i.buffer)
	i.pos += direction
	if i.pos < 0 {
		i.pos = 0
	}
	if i.pos >= max {
		i.pos = max
	}
}

func (i *inputBuffer) put(r rune) {
	i.buffer = append(i.buffer, r)
	copy(i.buffer[i.pos+1:], i.buffer[i.pos:])
	i.buffer[i.pos] = r
	i.pos++
}

func (i *inputBuffer) handle(ev *tcell.EventKey) (rerender, query_changed bool) {
	switch ev.Key() {
	case tcell.KeyBackspace2:
		fallthrough
	case tcell.KeyBackspace:
		b := i.backspace()
		return b, b
	case tcell.KeyDelete:
		b := i.delete()
		return b, b
	case tcell.KeyLeft:
		i.advance(-1)
		return true, false
	case tcell.KeyRight:
		i.advance(+1)
		return true, false
	case tcell.KeyRune:
		i.put(ev.Rune())
	default:
		return false, false
	}
	return true, true
}
