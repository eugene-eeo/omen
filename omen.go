package main

import "os"
import "fmt"
import "github.com/gdamore/tcell"

func die(what string) {
	fmt.Fprintf(os.Stderr, "rtpreview: %s", what)
	os.Exit(1)
}

func main() {
	opts := parseFlags()

	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	screen, err := tcell.NewScreen()
	if err != nil {
		die("cannot initialise tcell")
	}
	if err := screen.Init(); err != nil {
		die(err.Error())
	}
	screen.Clear()
	ib := inputBuffer{buffer: []rune{}}
	pm := newPreviewManager(opts)
	quit := make(chan bool)

	tcell_events := make(chan tcell.Event)
	go func() {
		for {
			tcell_events <- screen.PollEvent()
		}
	}()

	pm.listen()

	width := 100
	height := 25

	go func() {
		for {
			select {
			case pd := <-pm.sink:
				if pd.uid == pm.uid && 1+pd.lineNo < height {
					y := 1 + pd.lineNo
					unicodeCells([]rune(pd.line), width, true, func(x int, r rune, _ int) {
						screen.SetContent(x, y, r, nil, tcell.StyleDefault)
					})
					screen.Sync()
				}

			case ev := <-tcell_events:
				switch ev.(type) {
				case *tcell.EventResize:
					x := ev.(*tcell.EventResize)
					width, height = x.Size()
					pm.maxLines = height - 2
					pm.maxLineLength = width
					// Treat resizing as a new preview request
					// this shouldn't occur too frequently, but it's nice to handle it
					screen.Clear()
					pm.debouncePreview(string(ib.buffer))
					drawPrompt(screen, &ib, width)
					screen.Sync()

				case *tcell.EventKey:
					x := ev.(*tcell.EventKey)
					switch x.Key() {
					case tcell.KeyEscape:
						quit <- false
					case tcell.KeyEnter:
						quit <- true
					default:
						rerender, query_changed := ib.handle(x)
						if query_changed {
							screen.Clear()
							pm.debouncePreview(string(ib.buffer))
						}
						if rerender {
							drawPrompt(screen, &ib, width)
							screen.Sync()
						}
					}
				}
			}
		}
	}()

	drawPrompt(screen, &ib, width)
	screen.Sync()
	selected := <-quit
	screen.Fini()
	if selected {
		fmt.Println(string(ib.buffer))
	}
}
