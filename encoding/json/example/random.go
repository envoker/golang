package main

import (
	"bytes"
	"math/rand"
	"time"
)

var strList = []string{
	"привет",
	"мир",
	"природа",
	"работа",
	"учеба",
	"дружба",
	"смех",
	"жизнь",
	"родина",
	"возрадоваться",
	"успокоиться",
	"座高楼 - высокое здание;",
	"个很老的人 - очень старый человек;",
	"个非常好的朋友 - очень хороший друг;",
}

type InstanceRandomizer interface {
	InitRandomInstance(r *rand.Rand) (err error)
}

func newRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func randomBoolean(r *rand.Rand) bool {
	const n = 1000
	return (r.Intn(2*n) < n)
}

func randomString(r *rand.Rand) string {

	buffer := new(bytes.Buffer)

	wordCount := r.Intn(5)
	for i := 0; i < wordCount; i++ {

		if i > 0 {
			buffer.WriteRune(' ')
		}

		buffer.WriteString(strList[r.Intn(len(strList))])
	}

	return string(buffer.Bytes())
}
