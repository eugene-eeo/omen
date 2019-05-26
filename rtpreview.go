package main

import "os"
import "fmt"
import "github.com/gdamore/tcell"

func die(what string) {
	fmt.Fprintf(os.Stderr, "rtpreview: %s", what)
	os.Exit(1)
}

func main() {
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
	pm := previewManager{
		sink:           make(chan previewDone, 10),
		queue:          make(chan string),
		debouncedQueue: make(chan string),
	}
	quit := make(chan bool)

	tcell_events := make(chan tcell.Event)
	go func() {
		for {
			tcell_events <- screen.PollEvent()
		}
	}()

	preview_channel := pm.listen()

	width := 100
	height := 25

	go func() {
		for {
			select {
			case pd := <-preview_channel:
				if 1+pd.lineNo < height {
					y := 2 + pd.lineNo
					unicodeCells([]rune(pd.line), width-2, false, func(x int, r rune) {
						screen.SetContent(x, y, r, nil, tcell.StyleDefault)
					})
					screen.Sync()
				}

			case ev := <-tcell_events:
				switch ev.(type) {
				case *tcell.EventResize:
					x := ev.(*tcell.EventResize)
					width, height = x.Size()

				case *tcell.EventKey:
					x := ev.(*tcell.EventKey)
					switch x.Key() {
					case tcell.KeyEscape:
						quit <- false
					case tcell.KeyEnter:
						quit <- true
					default:
						changed := ib.handle(x)
						if !changed {
							continue
						}
						screen.Clear()
						drawPrompt(screen, &ib, width)
						screen.Sync()
						pm.debouncePreview(string(ib.buffer))
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