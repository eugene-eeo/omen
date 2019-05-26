package main

import "fmt"
import "time"
import "github.com/google/shlex"

type previewManager struct {
	uid           uint
	current       *preview
	sink          chan previewLine
	queue         chan string
	killChan      chan struct{}
	doneChan      chan struct{}
	maxLines      int
	maxLineLength int
	options       *cliOptions
}

func newPreviewManager(opt *cliOptions) *previewManager {
	return &previewManager{
		sink:          make(chan previewLine, 10),
		queue:         make(chan string, 5),
		killChan:      make(chan struct{}),
		doneChan:      make(chan struct{}),
		maxLines:      100,
		maxLineLength: 100,
		options:       opt,
	}
}

func (p *previewManager) listen() {
	go func() {
		interval := 50 * time.Millisecond
		timer := time.NewTimer(interval)
		query := ""
		for {
			select {
			case query = <-p.queue:
				timer.Reset(interval)
			case <-timer.C:
				p.perform(query)
			}
		}
	}()
}

func (p *previewManager) debouncePreview(query string) {
	p.queue <- query
}

func (p *previewManager) perform(query string) {
	if p.current != nil {
		p.current.kill()
		p.current = nil
	}

	parts, err := shlex.Split(fmt.Sprintf(p.options.cmdFormat, query))
	if err != nil && len(parts) > 0 {
		return
	}

	p.uid++
	p.current = &preview{
		uid:           p.uid,
		cmd:           parts[0],
		args:          parts[1:],
		killChan:      p.killChan,
		doneChan:      p.doneChan,
		sink:          p.sink,
		maxLines:      p.maxLines,
		maxLineLength: p.maxLineLength,
	}
	go p.current.start()
}
