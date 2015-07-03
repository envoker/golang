package main

import (
	"bytes"
	"math/rand"
	
	"github.com/envoker/golang/testing/random"
)

var randomWordTable = []string{
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
}

func randString(r *rand.Rand) string {

	var buffer *bytes.Buffer
	buffer = new(bytes.Buffer)

	maxLen := 5	
	countWords := random.Intn(r, maxLen)

	for i := 0; i < countWords; i++ {

		j := r.Intn(len(randomWordTable))

		buffer.WriteString(randomWordTable[j])
		buffer.WriteByte(' ')
	}

	return string(buffer.Bytes())
}
