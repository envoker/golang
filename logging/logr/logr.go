package logr

import (
	"errors"
	"fmt"
	"os"
)

type Writer struct {
	fileName   string
	fw         *fileWriter
	rotateChan chan bool
}

func New(fileName string, maxSize int64, countFiles int) (*Writer, error) {

	if maxSize < 0 {
		return nil, errors.New("maxSize < 0")
	}

	if (countFiles < 1) || (countFiles > 100) {
		return nil, errors.New("count backup files must be [1..100]")
	}

	fw := &fileWriter{maxSize: maxSize}

	err := fw.Open(fileNameWithSuffix(fileName, -1))
	if err != nil {
		return nil, err
	}

	rotateChan := make(chan bool)

	go loopRotate(rotateChan, fileName, countFiles)

	return &Writer{
		fileName:   fileName,
		fw:         fw,
		rotateChan: rotateChan,
	}, nil
}

func (w *Writer) Close() error {

	err := w.fw.Close()
	w.fw = nil
	if err != nil {
		return err
	}

	w.rotateChan <- false
	<-w.rotateChan

	return nil
}

func (w *Writer) Write(data []byte) (n int, err error) {

	if w.fw.available() < int64(len(data)) {

		w.fw.Close()

		w.rotateChan <- true
		<-w.rotateChan // wait to rename file

		err = w.fw.Open(fileNameWithSuffix(w.fileName, -1))
		if err != nil {
			return 0, err
		}
	}

	n, err = w.fw.Write(data)
	return
}

func loopRotate(rotateChan chan bool, fileName string, countFiles int) {
	for {
		if <-rotateChan {
			oldpath := fileNameWithSuffix(fileName, -1)
			if fileExist(oldpath) {
				newpath := fileNameWithSuffix(fileName, 0)
				os.Rename(oldpath, newpath)
			}
			rotateChan <- true
			rotate(fileName, countFiles)
		} else {
			rotateChan <- false
			return
		}
	}
}

func rotate(fileName string, countFiles int) {

	index := countFiles

	nameLast := fileNameWithSuffix(fileName, index)
	if fileExist(nameLast) {
		os.Remove(nameLast)
	}

	for index > 0 {
		index--
		oldpath := fileNameWithSuffix(fileName, index)
		if fileExist(oldpath) {
			newpath := fileNameWithSuffix(fileName, index+1)
			os.Rename(oldpath, newpath)
		}
	}
}

func fileNameWithSuffix(fileName string, index int) string {
	if index < 0 {
		return fmt.Sprintf("%s.log", fileName)
	}
	return fmt.Sprintf("%s.%d.log", fileName, index)
}

func fileExist(fileName string) bool {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	}
	return true
}
