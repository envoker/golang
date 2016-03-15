package logr

import (
	"os"
)

type fileWriter struct {
	file    *os.File
	size    int64
	maxSize int64
}

func (fw *fileWriter) Open(fileName string) error {

	if err := fw.Close(); err != nil {
		return nil
	}

	var size int64

	if fi, err := os.Stat(fileName); err == nil {
		size = fi.Size()
	}

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	fw.file = file
	fw.size = size

	return nil
}

func (fw *fileWriter) Close() error {
	if fw.file != nil {
		err := fw.file.Close()
		fw.file = nil
		return err
	}
	return nil
}

func (fw *fileWriter) Write(data []byte) (n int, err error) {
	n, err = fw.file.Write(data)
	fw.size += int64(n)
	return
}

func (fw *fileWriter) available() int64 {
	return fw.maxSize - fw.size
}
