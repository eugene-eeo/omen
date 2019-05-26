package main

import "time"

type previewManager struct {
	uid      uint
	current  *preview
	sink     chan previewLine
	queue    chan string
	killChan chan struct{}
	doneChan chan struct{}
	maxLines int
}

func newPreviewManager() *previewManager {
	return &previewManager{
		sink:     make(chan previewLine, 10),
		queue:    make(chan string),
		killChan: make(chan struct{}),
		doneChan: make(chan struct{}),
		maxLines: 100,
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
	}
	p.uid++
	p.current = &preview{
		uid:      p.uid,
		cmd:      "ag",
		args:     []string{"--", query},
		killChan: p.killChan,
		doneChan: p.doneChan,
		sink:     p.sink,
		lines:    p.maxLines,
	}
	go p.current.start()
}
