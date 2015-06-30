package main

import (
	"bytes"
	"math/rand"
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

func RandomRange(r *rand.Rand, min int, max int) int {
	return (min + r.Intn(max-min+1))
}

func RandomString(r *rand.Rand) string {

	var buffer *bytes.Buffer
	buffer = new(bytes.Buffer)

	minLen := 2
	maxLen := 5
	countWords := RandomRange(r, minLen, maxLen)

	for i := 0; i < countWords; i++ {

		j := r.Intn(len(randomWordTable))

		buffer.WriteString(randomWordTable[j])
		buffer.WriteByte(' ')
	}

	return string(buffer.Bytes())
}
