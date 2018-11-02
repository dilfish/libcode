package tools

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
    "sync"
	"time"
)

type AppendStruct struct {
	file   *os.File
	c      chan os.Signal
	cClose chan bool
	fn     string
	err    error
    pid    int
    lock   sync.Mutex
}

func openFile(fn string) (*os.File, error) {
	return os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
}

func (as *AppendStruct) wait() {
	for {
		select {
		case <-as.c:
			io.WriteString(os.Stderr, "we got an signal, restart at "+TimeStr())
			f, err := openFile(as.fn)
			if err != nil {
				as.err = err
				log.Println("open file:", err)
				time.Sleep(time.Second)
				continue
			}
            as.lock.Lock()
            if as.file != nil {
                as.file.Close()
            }
			as.file = f
            as.lock.Unlock()
		case <-as.cClose:
			signal.Reset(syscall.SIGUSR1)
			return
		}
	}
}

func NewAppender(fn string) (*AppendStruct, error) {
	var as AppendStruct
	f, err := openFile(fn)
	if err != nil {
		return nil, err
	}
	as.file = f
	as.fn = fn
	as.c = make(chan os.Signal)
	as.cClose = make(chan bool)
    as.pid = os.Getpid()
    signal.Notify(as.c, syscall.SIGUSR1)
	go as.wait()
	return &as, nil
}

func (as *AppendStruct) Close() {
    as.lock.Lock()
    defer as.lock.Unlock()
	as.cClose <- true
	as.file.Close()
    as.file = nil
	close(as.c)
	close(as.cClose)
}

func (as *AppendStruct) Write(bt []byte) (int, error) {
	if as.err != nil {
		return 0, as.err
	}
    as.lock.Lock()
    defer as.lock.Unlock()
	n, err := as.file.Write(bt)
	if err != nil {
		as.err = err
	}
	return n, err
}

func Daemon() {
	os.Stdout.Close()
	os.Stdin.Close()
	os.Stdout = nil
	os.Stdin = nil
}

func InitLog(fn, prefix string) *log.Logger {
	as, err := NewAppender(fn)
	if err != nil {
		return nil
	}
	if prefix == "" {
		prefix = "DefAppendLogger "
	}
	if prefix[len(prefix)-1] != ' ' {
		prefix = prefix + " "
	}
	return log.New(as, prefix, log.LstdFlags|log.Lshortfile)
}
