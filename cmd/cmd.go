package cmd

import (
	"bytes"
	"context"
	"github.com/google/uuid"
	"io"
	"log"
	"os/exec"
	"time"
)

var (
	// 默认超时时间2分钟
	defaultTimeout = 2 * 60 * time.Second
)

type Work struct {
	Command string
	Args    []string
	TimeOut time.Duration
	id      string
	cmd     *exec.Cmd
	context context.Context
	Out     string
	Error   string
	IsError bool
}

func (w *Work) Run() {
	w.id = uuid.New().String()
	if w.TimeOut == 0 {
		w.TimeOut = defaultTimeout
	}
	ctxt, cancel := context.WithTimeout(context.Background(), w.TimeOut)
	defer cancel()
	w.context = ctxt
	cmd := exec.CommandContext(ctxt, w.Command, w.Args...)
	w.cmd = cmd

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer

	//var errStdout, errStderr error
	stdout := io.MultiWriter(&outBuf)
	stderr := io.MultiWriter(&errBuf)

	log.Printf("ID: %s, Content: %s, Timeout: %v s", w.id, w.cmd.String(), w.TimeOut.Seconds())
	err := cmd.Start()
	if err != nil {
		log.Printf("ID: %s, start has error: %s", w.id, err)
		return
	}
	go func() {
		_, _ = io.Copy(stdout, stdoutIn)
	}()
	go func() {
		_, _ = io.Copy(stderr, stderrIn)
	}()
	log.Printf("ID: %s, PID: %d", w.id, w.cmd.Process.Pid)
	err = cmd.Wait()
	if err != nil {
		w.IsError = true
		w.Out = outBuf.String()
		w.Error = errBuf.String()
		cErr := ctxt.Err()
		if cErr != nil {
			w.Error = "timeout"
			err = cErr
		}
		log.Printf("ID: %s, PID: %d, exec has error: %s, error output: %s", w.id, w.cmd.Process.Pid, err, w.Error)
	} else {
		w.IsError = false
		w.Out = outBuf.String()
		w.Error = errBuf.String()
		log.Printf("ID: %s, PID: %d, exec end and output: %s, error output: %s", w.id, w.cmd.Process.Pid, w.Out, w.Error)
	}
}

func (w *Work) Stop() {
	log.Printf("ID: %s, PID: %d, stoping", w.id, w.cmd.Process.Pid)
	err := w.cmd.Process.Kill()
	if err != nil {
		log.Printf("ID: %s, stop error, error output: %s", w.id, err)
	} else {
		log.Printf("ID: %s, stop success", w.id)
	}
}
