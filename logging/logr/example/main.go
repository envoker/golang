package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/envoker/golang/logging/logr"
)

const (
	kilobyte  = 1024
	megabyte  = 1024 * kilobyte
	gigabyte  = 1024 * megabyte
	terabyte  = 1024 * gigabyte
	petabyte  = 1024 * terabyte
	exabyte   = 1024 * petabyte
	zettabyte = 1024 * exabyte
	yottabyte = 1024 * zettabyte
)

func main() {
	logExample()
}

func test1() {

	wg := new(sync.WaitGroup)
	ch := make(chan int, 20)

	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case k, ok := <-ch:
				if !ok {
					return
				}
				fmt.Println(k)
			}
		}
	}()

	for i := 0; i < 167; i++ {
		ch <- i
	}

	close(ch)

	wg.Wait()
}

func logSimple() {

	w := logr.New(logr.Config{
		FileName: "./test/mylog",
		FileSize: 512 * kilobyte,
		Count:    5,
	})
	defer w.Close()

	logger := log.New(w, "info ", log.Ldate|log.Lmicroseconds)

	logger.Println("out of range 1")

	w.Close()

	logger.Println("out of range 2")
}

func zipExample() {

	buffer := new(bytes.Buffer)

	w := zip.NewWriter(buffer)

	fileName := "./test/mylog"
	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	f, err := w.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write(fileBytes)
	if err != nil {
		log.Fatal(err)
	}

	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(fileName + ".zip")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(buffer.Bytes())
	if err != nil {
		log.Fatal(err)
	}
}

func gzipExample() error {

	fileName := "./test/mylog.0"

	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf("%s.gz", fileName))
	if err != nil {
		return err
	}
	defer file.Close()

	bw := bufio.NewWriter(file)
	defer bw.Flush()

	archive := gzip.NewWriter(bw)
	defer archive.Close()

	if _, err = archive.Write(fileBytes); err != nil {
		return err
	}

	return nil
}

func logExample() {

	dirName := "test"

	os.Mkdir(dirName, os.ModePerm)

	w := logr.New(logr.Config{
		FileName: filepath.Join(dirName, "mylog"),
		FileSize: 512 * kilobyte,
		Count:    5,
	})
	defer w.Close()

	logger := log.New(w, "info ", log.Ldate|log.Lmicroseconds)

	const n = 1000
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(id int) {
			for j := 0; j < 1000; j++ {
				logger.Printf("goroutine: %d; cmd: %d\n", id, j)
				runtime.Gosched()
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
