package main

import "time"
import "sync"
import "testing"

func getPreview() (*preview, chan previewLine) {
	sink := make(chan previewLine)
	p := &preview{
		uid:      0,
		cmd:      "yes",
		args:     []string{},
		buffer:   make([]byte, 0, 10),
		killChan: make(chan struct{}),
		sink:     sink,
		maxLines: 100,
		wg:       &sync.WaitGroup{},
	}
	return p, sink
}

func TestPreview(t *testing.T) {
	// Check that preview sends the correct # of lines with correct lines
	timer := time.NewTimer(time.Second)
	done := make(chan bool)
	p, sink := getPreview()

	go func() {
		for i := 0; i < 100; i++ {
			x := <-sink
			if !(x.uid == 0 && string(x.line) == "y") {
				t.Errorf("wrong output: expected y, got '%s' instead.", string(x.line))
			}
		}
		done <- true
	}()

	p.start()
	select {
	case <-done:
	case <-timer.C:
		t.Fatal("test didn't finish in time")
	}
}

func TestPreviewKill(t *testing.T) {
	// Check that we can actually kill previews
	p, sink := getPreview()
	doneTimer := time.NewTimer(time.Second)
	done := make(chan bool)

	go func() {
		for _ = range sink {
		}
		done <- true
	}()

	go func() {
		time.Sleep(time.Millisecond * 50)
		p.kill()
		close(sink)
	}()

	p.start()
	select {
	case <-done:
	case <-doneTimer.C:
		t.Fatal("test didn't finish in time")
	}
}
