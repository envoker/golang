package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/envoker/golang/logging/logr"
)

func main() {

	const (
		kilobyte = 1024
		megabyte = 1024 * kilobyte
	)

	dir := "logs"
	os.Mkdir(dir, os.ModePerm)

	w, err := logr.New(filepath.Join(dir, "test"), 100*kilobyte, 7)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	logger := log.New(w, "tag ", log.Ldate|log.Lmicroseconds)

	r := newRand()
	for i := 0; i < 100000; i++ {
		logger.Println(i, randString(r))
	}
}
