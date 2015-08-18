package daylog

import (
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/envoker/golang/time/date"
)

func removeOld(dirName string, daysNumber int) error {

	fis, err := ioutil.ReadDir(dirName)
	if err != nil {
		return err
	}

	dateNow := date.Now()

	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		fileName := fi.Name()
		dateFile, err := dateFromFileName(fileName)
		if err != nil {
			continue
		}

		if dateNow.Sub(dateFile) >= daysNumber {
			os.Remove(fileName)
		}
	}

	return nil
}

func removeOldWorker(quit chan struct{}, wg *sync.WaitGroup, dirName string, daysNumber int) {

	defer wg.Done()

	ticker := time.NewTicker(2 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			return

		case <-ticker.C:
			removeOld(dirName, daysNumber)
		}
	}
}
