package daylog

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/envoker/golang/time/date"
)

type rotator struct {
	dirName    string
	daysNumber int
}

func (r *rotator) Rotate() error {
	return removeOld(r.dirName, r.daysNumber)
}

func removeOld(dirName string, daysNumber int) error {

	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		return err
	}

	dateNow := date.Now()

	for _, fileInfo := range files {
		if !fileInfo.IsDir() {
			fileName := fileInfo.Name()
			dateFile, err := dateFromFileName(fileName)
			if err == nil {
				if dateNow.Sub(dateFile) >= daysNumber {
					if err := os.Remove(filepath.Join(dirName, fileName)); err != nil {
						log.Println(err)
					}
				}
			}
		}
	}

	return nil
}
