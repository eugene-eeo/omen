package main

import "os/exec"
import "bufio"

type preview struct {
	uid      uint
	cmd      string
	args     []string
	killChan chan struct{}
	doneChan chan struct{}
	sink     chan previewLine
	lines    int
}

type previewLine struct {
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
			n++
			p.sink <- previewLine{p.uid, n, scanner.Text()}
		}
	}()

	go func() {
		// When we receive a kill-signal, just kill the process and send
		// a signal to p.done.
		<-p.killChan
		cmd.Process.Kill()
		cmd.Process.Release()
		p.doneChan <- struct{}{}
	}()
}

func (p *preview) kill() {
	p.killChan <- struct{}{}
	<-p.doneChan
}
