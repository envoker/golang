package main

import (
	"log"
	"math/rand"
	"runtime"
	"sync"

	"github.com/envoker/golang/logging/daylog"
	"github.com/envoker/golang/logging/logl"
)

func main() {
	if err := loglFiller(); err != nil {
		log.Println(err)
	}
}

func logFiller() error {

	w, err := daylog.New("./test", 10, "utc ")
	if err != nil {
		return err
	}
	defer w.Close()

	logger := log.New(w, "ios ", 0)

	const n = 100
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go testLoop(wg, logger, i)
	}
	wg.Wait()

	return nil
}

func testLoop(wg *sync.WaitGroup, logger *log.Logger, index int) {

	defer wg.Done()

	r := rand.New(rand.NewSource(int64(index)))
	for i := 0; i < 100; i++ {
		logger.Printf("routine(%d):%s", index, randString(r, 40))
		runtime.Gosched()
	}
}

func loglFiller() error {

	w, err := daylog.New("./test", 10, "test ")
	if err != nil {
		return err
	}
	defer w.Close()

	logger := logl.New(w, "ios ", logl.LEVEL_DEBUG, 0)

	const n = 100
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go testLoopLogl(wg, logger, i)
	}
	wg.Wait()

	return nil
}

func testLoopLogl(wg *sync.WaitGroup, logger *logl.Logger, index int) {

	defer wg.Done()

	r := rand.New(rand.NewSource(int64(index)))
	for i := 0; i < 100; i++ {
		logger.Infof("routine(%d):%s", index, randString(r, 40))
		runtime.Gosched()
	}
}
