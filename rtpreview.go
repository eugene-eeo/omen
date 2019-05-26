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
	quit := make(chan struct{})

	tcell_events := make(chan tcell.Event)
	go func() {
		for {
			tcell_events <- screen.PollEvent()
		}
	}()

	preview_channel := pm.listen()

	go func() {
		for {
			select {
			case pd := <-preview_channel:
				for dx, r := range []rune(pd.line) {
					screen.SetCell(dx, 2+pd.lineNo, tcell.StyleDefault, r)
				}
				screen.Sync()

			case ev := <-tcell_events:
				switch ev.(type) {
				case *tcell.EventKey:
					x := ev.(*tcell.EventKey)
					switch x.Key() {
					case tcell.KeyEscape:
						quit <- struct{}{}
					default:
						changed := ib.handle(x)
						if !changed {
							continue
						}
						screen.Clear()
						drawPrompt(screen, &ib)
						screen.Sync()
						pm.debouncePreview(string(ib.buffer))
					}
				}
			}
		}
	}()

	drawPrompt(screen, &ib)
	screen.Sync()
	<-quit
	screen.Fini()
}
