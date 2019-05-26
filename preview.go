package main

import "os/exec"
import "bufio"

type preview struct {
	uid   uint
	cmd   string
	args  []string
	kill  chan struct{}
	done  chan struct{}
	sink  chan previewDone
	lines int
}

type previewDone struct {
	uid    uint
	lineNo int
	line   string
}

func (p *preview) start() {
	cmd := exec.Command(p.cmd, p.args...)
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	// 1 goroutine for listening for output
	// 1 goroutine for waiting for kill signal (if any)

	go func() {
		n := 0
		scanner := bufio.NewScanner(stdout)
		for n < p.lines && scanner.Scan() {
			p.sink <- previewDone{p.uid, n, scanner.Text()}
			n++
		}
	}()

	go func() {
		// When we're done or we receive a kill-signal, just kill the damn process,
		// and send a signal to p.done.
		<-p.kill
		cmd.Process.Kill()
		cmd.Process.Release()
		p.done <- struct{}{}
	}()
}

func (p *preview) destroy() {
	p.kill <- struct{}{}
	<-p.done
}
