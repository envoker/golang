package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"

	"github.com/envoker/golang/logging/daylog"
	"github.com/envoker/golang/logging/logl"
)

func main() {

	fn := logFiller
	//fn := loglFiller

	if err := fn(); err != nil {
		log.Println(err)
	}
}

func logFiller() error {

	w, err := daylog.New("./logs", 10, "utc ")
	if err != nil {
		return err
	}
	defer w.Close()

	logger := log.New(w, "ios ", 0)

	const n = 100
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go logLoop(wg, logger, i)
	}
	wg.Wait()

	return nil
}

func logLoop(wg *sync.WaitGroup, logger *log.Logger, index int) {

	defer wg.Done()

	r := rand.New(rand.NewSource(int64(index)))
	for i := 0; i < 100; i++ {
		logger.Printf("routine(%d): %s", index, randString(r, 15))
		runtime.Gosched()
	}
}

func loglFiller() error {

	w, err := daylog.New("./test", 10, "test ")
	if err != nil {
		return err
	}
	defer w.Close()

	logger := logl.New(w, "ios ", logl.LEVEL_WARNING, 0)

	const n = 100
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go loglLoop(wg, logger, i)
	}
	wg.Wait()

	return nil
}

func loglLoop(wg *sync.WaitGroup, logger *logl.Logger, index int) {

	defer wg.Done()

	r := rand.New(rand.NewSource(int64(index)))
	for i := 0; i < 100; i++ {

		m := fmt.Sprintf("routine(%d): %s", index, randString(r, 15))

		switch k := r.Intn(6); k {
		case 0:
			logger.Fatal(m)
		case 1:
			logger.Error(m)
		case 2:
			logger.Warning(m)
		case 3:
			logger.Info(m)
		case 4:
			logger.Debug(m)
		case 5:
			logger.Trace(m)
		}

		runtime.Gosched()
	}
}
