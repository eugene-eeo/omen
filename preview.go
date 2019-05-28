package main

import "bytes"
import "os/exec"
import "bufio"
import "sync"

type preview struct {
	uid      uint
	cmd      string
	args     []string
	buffer   []byte
	killChan chan struct{}
	sink     chan previewLine
	maxLines int
	wg       *sync.WaitGroup
}

type previewLine struct {
	uid    uint
	lineNo int
	line   []rune
}

func (p *preview) start() {
	cmd := exec.Command(p.cmd, p.args...)
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	// 1 goroutine for listening for output
	// 1 goroutine for waiting for kill signal (if any)
	p.wg.Add(2)

	go func() {
		n := 0
		scanner := bufio.NewScanner(stdout)
		scanner.Buffer(p.buffer, 0)
		for n < p.maxLines && scanner.Scan() {
			n++
			p.sink <- previewLine{p.uid, n, bytes.Runes(scanner.Bytes())}
		}
		stdout.Close()
		p.wg.Done()
	}()

	go func() {
		// When we receive a kill-signal, just kill the process and send
		// a signal to p.done.
		<-p.killChan
		if cmd.Process != nil {
			cmd.Process.Kill()
			cmd.Process.Release()
		}
		p.wg.Done()
	}()
}

func (p *preview) kill() {
	p.killChan <- struct{}{}
	p.wg.Wait()
}
