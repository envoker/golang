package main

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/envoker/golang/logging/clog"
)

func logTest() {

	logger, err := clog.NewLogger(clog.Config{
		Dir:          "./log",
		CountRecords: 1000,
		Level:        clog.LEVEL_INFO,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer logger.Close()

	var n = 100
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go procLogTest(wg, i, logger)
	}

	logger.Error(fmt.Sprintf("%d %s", 213, "Конец"))
	logger.Warning(":)")
	logger.Info("information")
	logger.Debug("Конец")

	wg.Wait()
}

func procLogTest(wg *sync.WaitGroup, id int, logger *clog.Logger) {

	defer wg.Done()

	var level clog.Level
	r := rand.New(rand.NewSource(int64(id)))
	for i := 0; i < 500; i++ {
		level = clog.Level(1 + r.Intn(6))
		logger.Log(level, "thread: %03d | %s", id, randString(r))
	}
}
