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
}

func useLogger(logger *logl.Logger) {

	logger.Fatal("fatal message")
	logger.Error("error message")
	logger.Warning("warning message")
	logger.Info("info message")
	logger.Debug("debug message")
	logger.Trace("trace message")

	logger.Fatalf("fatal message: %d", 1)
	logger.Errorf("error message: %d", 2)
	logger.Warningf("warning message: %d", 3)
	logger.Infof("info message: %d", 4)
	logger.Debugf("debug message: %d", 5)
	logger.Tracef("trace message: %d", 6)
}

func exampleLogStdout() {

	logger := logl.New(os.Stdout, logl.LEVEL_DEBUG, logl.Ltime)

	useLogger(logger)
}

func exampleLogOff() {

	logger := logl.New(os.Stdout, 0, logl.Ltime)

	useLogger(logger)
}

func exampleLogFile() error {

	file, err := os.OpenFile("./test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	defer w.Flush()

	logger := logl.New(w, logl.LEVEL_DEBUG, logl.Ldate|logl.Lmicroseconds)

	useLogger(logger)

	return nil
}
