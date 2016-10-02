package daylog

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/toelsiba/date"
)

type flusher interface {
	flush() error
}

type rotator struct {
	dir        string
	daysNumber int
}

func (r *rotator) Rotate() {
	if err := removeOld(r.dir, r.daysNumber); err != nil {
		log.Println(err)
	}
}

func removeOld(dir string, daysNumber int) error {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	dateNow := date.CurrentDate()

	for _, fileInfo := range files {
		if !fileInfo.IsDir() {
			fileName := fileInfo.Name()
			dateFile, err := dateFromFileName(fileName)
			if err == nil {
				if dateFile.DaysTo(dateNow) >= daysNumber {
					if err := os.Remove(filepath.Join(dir, fileName)); err != nil {
						log.Println(err)
					}
				}
			}
		}
	}

	return nil
}

func worker(quit chan bool, f flusher, r rotator) {

	r.Rotate()

	flushTicker := time.NewTicker(5 * time.Minute)
	defer flushTicker.Stop()

	rotateTicker := time.NewTicker(12 * time.Hour)
	defer rotateTicker.Stop()

	for {
		select {
		case <-quit:
			quit <- true
			return

		case <-flushTicker.C:
			f.flush()

		case <-rotateTicker.C:
			r.Rotate()
		}
	}
}
