package daylog

import (
	"log"
	"time"
)

type flusher interface {
	Flush() error
}

func worker(quit chan bool, f flusher, r rotator) {

	flushTicker := time.NewTicker(12 * time.Hour)
	defer flushTicker.Stop()

	rotateTicker := time.NewTicker(5 * time.Minute)
	defer rotateTicker.Stop()

	if err := r.Rotate(); err != nil {
		log.Println(err)
	}

	for {
		select {
		case <-quit:
			quit <- true
			return

		case <-flushTicker.C:
			if err := f.Flush(); err != nil {
				log.Println(err)
			}

		case <-rotateTicker.C:
			if err := r.Rotate(); err != nil {
				log.Println(err)
			}
		}
	}
}
