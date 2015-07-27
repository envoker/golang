package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"log/syslog"
	"os"
	"sync"
	"time"

	"github.com/envoker/golang/logging/nalog"
)

func main() {

	// logFiller
	// logParser()
	// logGolang

	if err := logFiller(); err != nil {
		fmt.Println(err)
	}
}

func logFiller() error {

	l, err := nalog.New("./", nalog.LEVEL_WARNING)
	if err != nil {
		return err
	}
	defer l.Close()

	logger := l.Logger()

	wg := new(sync.WaitGroup)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go testFn(wg, logger, i)
	}

	wg.Wait()

	return nil
}

func testFn(wg *sync.WaitGroup, logger nalog.Logger, index int) {

	defer wg.Done()

	r := newRand()
	for i := 0; i < 10; i++ {
		logger.Logf(
			nalog.Level(r.Intn(4)+1),
			"routine(%d):%s",
			index, randString(r, 20),
		)
		time.Sleep(300 * time.Millisecond)
	}
}

func logParser() error {

	fileName := "./2015-07-24.log"

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	var (
		r = bufio.NewReader(file)
		d = xml.NewDecoder(r)
	)

	var record struct {
		XMLName xml.Name `xml:"msg"`
		Time    string   `xml:"time,attr"`
		Level   string   `xml:"level,attr"`
		Message string   `xml:",innerxml"`
	}

	for {
		if err := d.Decode(&record); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		fmt.Println(record.Message)
	}

	return nil
}

func logGolang() error {

	l3, err := syslog.New(syslog.LOG_ERR, "GoExample")
	//l3, err := syslog.Dial("udp", "localhost", syslog.LOG_ERR, "GoExample") // connection to a log daemon
	defer l3.Close()
	if err != nil {
		log.Fatal("error")
	}

	//l3.Emerg("emergency")
	l3.Alert("alert")
	l3.Crit("critical")
	l3.Err("error")
	l3.Warning("warning")
	l3.Notice("notice")
	l3.Info("information")
	l3.Debug("debug")
	l3.Write([]byte("write"))

	return nil
}
