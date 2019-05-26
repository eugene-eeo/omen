package main

import "sync"
import "time"

type previewManager struct {
	uid           uint
	current       *preview
	sink          chan previewLine
	queue         chan string
	killChan      chan struct{}
	maxLines      int
	maxLineLength int
	buffer        []byte
	options       *cliOptions
	wg            sync.WaitGroup
}

func newPreviewManager(opt *cliOptions) *previewManager {
	return &previewManager{
		sink:     make(chan previewLine, 5),
		queue:    make(chan string, 1),
		killChan: make(chan struct{}),
		options:  opt,
	}
}

func (p *previewManager) listen() {
	go func() {
		timer := time.NewTimer(p.options.debounceTime)
		query := ""
		for {
			select {
			case query = <-p.queue:
				timer.Reset(p.options.debounceTime)
			case <-timer.C:
				if query != "" || p.options.allowEmpty {
					p.perform(query)
				}
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

	parts, err := expandCommand(p.options.cmdFormat, query)
	if err != nil {
		return
	}

	if cap(p.buffer) < p.maxLineLength {
		p.buffer = make([]byte, 0, p.maxLineLength)
	}

	p.uid++
	p.current = &preview{
		uid:      p.uid,
		cmd:      parts[0],
		args:     parts[1:],
		buffer:   p.buffer,
		killChan: p.killChan,
		sink:     p.sink,
		maxLines: p.maxLines,
		wg:       &p.wg,
	}
	p.current.start()
}
