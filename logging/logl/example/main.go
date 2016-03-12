package main

import (
	"bufio"
	"os"

	"github.com/envoker/golang/logging/logl"
)

func main() {
	exampleLogStdout()
	exampleLogOff()
	exampleLogFile()
	exampleLogAscii()
}

func useLogger(l *logl.Logger) {

	l.Fatal("fatal message")
	l.Error("error message")
	l.Warning("warning message")
	l.Info("info message")
	l.Debug("debug message")
	l.Trace("trace message")

	l.Fatalf("fatal message: %d", 1)
	l.Errorf("error message: %d", 2)
	l.Warningf("warning message: %d", 3)
	l.Infof("info message: %d", 4)
	l.Debugf("debug message: %d", 5)
	l.Tracef("trace message: %d", 6)
}

func exampleLogStdout() {
	l := logl.New(os.Stdout, "", logl.LEVEL_DEBUG, logl.Ltime)
	useLogger(l)
}

func exampleLogOff() {
	l := logl.New(os.Stdout, "", 0, logl.Ltime)
	useLogger(l)
}

func exampleLogFile() {

	file, err := os.OpenFile("test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	defer w.Flush()

	logger := logl.New(w, "prefix ", logl.LEVEL_DEBUG, logl.Ldate|logl.Lmicroseconds)

	useLogger(logger)
}

func exampleLogAscii() {

	logger := logl.New(os.Stdout, "test ", logl.LEVEL_DEBUG, logl.Ldate|logl.Lmicroseconds)

	data := make([]byte, 128)
	for i := range data {
		data[i] = byte(i)
	}

	logger.Debug(string(data))
}
