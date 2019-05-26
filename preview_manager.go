package main

import "time"

type previewManager struct {
	uid            uint
	current        *preview
	sink           chan previewDone
	queue          chan string
	debouncedQueue chan string
}

func (p *previewManager) listen() chan previewDone {
	sink := make(chan previewDone)
	query := ""
	go func() {
		interval := 50 * time.Millisecond
		timer := time.NewTimer(interval)
		for {
			select {
			case query = <-p.queue:
				timer.Reset(interval)
			case <-timer.C:
				p.debouncedQueue <- query
			}
		}
	}()
	go func() {
		for {
			select {
			case request := <-p.debouncedQueue:
				p.perform(request)
			case pd := <-p.sink:
				if pd.uid == p.uid {
					sink <- pd
				}
			}
		}
	}()
	return sink
}

func (p *previewManager) debouncePreview(query string) {
	p.queue <- query
}

func (p *previewManager) perform(query string) {
	p.uid++
	if p.current != nil {
		p.current.destroy()
	}

	p.current = &preview{
		uid:   p.uid,
		cmd:   "ag",
		args:  []string{"--", query},
		kill:  make(chan struct{}, 1),
		done:  make(chan struct{}, 1),
		sink:  p.sink,
		lines: 100,
	}
	p.current.start()
}
