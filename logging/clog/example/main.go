package main

import (
	"encoding/json"
	"fmt"

	"github.com/envoker/golang/logging/clog"
)

func main() {

	//logTest()
	LogJsonConvertTest()
}

func LogJsonConvertTest() {

	var c, p clog.Config

	c.Dir = "./log"
	c.CountRecords = 17078875
	c.Level = clog.LEVEL_DEBUG

	fmt.Printf("%+v\n", c)

	b, err := json.MarshalIndent(&c, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))

	err = json.Unmarshal(b, &p)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", p)
}
