package random

import (
	"bytes"
	"math/rand"
)

type alphabet struct {
	lower []rune
	upper []rune
}

var (
	abcEnglish = alphabet{
		lower: []rune("abcdefghijklmnopqrstuvwxyz"),
		upper: []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
	}

	abcRussian = alphabet{
		lower: []rune("абвгдеёжзийклмнопрстуфхчцшщъыьэюя"),
		upper: []rune("АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"),
	}

	abcUkrainian = alphabet{
		lower: []rune("абвгґдеєжзиіїйклмнопрстуфхцчшщьюя"),
		upper: []rune("АБВГҐДЕЄЖЗИІЇЙКЛМНОПРСТУФХЦЧШЩЬЮЯ"),
	}
)

var (
	digits       = []rune("0123456789")
	specialRunes = []rune("! @ # $ % ^ & * ( ) _ + , . / ?")
	hieroglyphs  = []rune("座高楼个很老的人个非常好的朋友")
)

var randRunesSamples = [][]rune{
	digits,
	specialRunes,
	abcEnglish.lower,
	abcEnglish.upper,
	abcRussian.lower,
	abcRussian.upper,
	abcUkrainian.lower,
	abcUkrainian.upper,
	hieroglyphs,
}

func String(r *rand.Rand, maxLen int) string {

	if maxLen < 0 {
		maxLen = 0
	}

	var buffer bytes.Buffer

	n := Intn(r, maxLen)

	for i := 0; i < n; i++ {

		j := r.Intn(len(randRunesSamples))
		k := r.Intn(len(randRunesSamples[j]))

		buffer.WriteRune(randRunesSamples[j][k])
	}

	return buffer.String()
}
