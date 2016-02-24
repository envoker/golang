package logr

import (
	"errors"
	"fmt"
	"os"
)

func rotate(fileName string, fileCount int) error {

	if fileCount < 1 {
		return errors.New("fileCount < 1")
	}

	index := fileCount - 1

	nameLast := fileNameIndex(fileName, index)
	if fileExist(nameLast) {
		if err := os.Remove(nameLast); err != nil {
			return err
		}
	}

	for index > 0 {
		index--

		oldpath := fileNameIndex(fileName, index)
		if fileExist(oldpath) {
			newpath := fileNameIndex(fileName, index+1)
			if err := os.Rename(oldpath, newpath); err != nil {
				return err
			}
		}
	}

	if fileExist(fileName) {
		newpath := fileNameIndex(fileName, 0)
		if err := os.Rename(fileName, newpath); err != nil {
			return err
		}
	}
	/*
		file, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer file.Close()
	*/
	return nil
}

func fileNameIndex(fileName string, i int) string {
	return fmt.Sprintf("%s.%d", fileName, i)
}

func fileExist(fileName string) bool {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	}
	return true
}
