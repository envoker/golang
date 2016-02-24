package logr

import (
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

var (
	_ io.Writer = &Writer{}
	_ io.Closer = &Writer{}
)

type Config struct {
	FileName string
	FileSize int64
	Count    int
}

type Writer struct {
	config Config
	wg     *sync.WaitGroup
	ch     chan []byte
	closed *int32
}

func New(config Config) *Writer {

	w := &Writer{
		config: config,
		wg:     new(sync.WaitGroup),
		ch:     make(chan []byte, 20),
		closed: new(int32),
	}

	w.wg.Add(1)
	go writeLoop(w.wg, w.ch, config)

	return w
}

func (w *Writer) Close() error {

	if !atomic.CompareAndSwapInt32(w.closed, 0, 1) {
		return errors.New("log writer is closed")
	}

	close(w.ch)
	w.wg.Wait()

	return nil
}

func (w *Writer) Write(p []byte) (n int, err error) {

	if atomic.LoadInt32(w.closed) != 0 {
		return 0, errors.New("log writer is closed")
	}

	w.ch <- duplicate(p)

	return
}

func writeLoop(wg *sync.WaitGroup, ch <-chan []byte, config Config) {

	defer wg.Done()

	fw := &fileWriter{config: config}
	defer fw.Close()

	for {
		select {
		case data, ok := <-ch:
			if !ok {
				return
			}
			fw.Write(data)

		case <-time.After(time.Minute):
			fw.Flush()
		}
	}
}
