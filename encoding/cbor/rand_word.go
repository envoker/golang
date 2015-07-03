package cbor

import (
	"math/rand"

	"github.com/envoker/golang/testing/random"
)

type alphabet struct {
	Lower []rune
	Upper []rune
}

var (
	abcEnglish = alphabet{
		Lower: []rune("abcdefghijklmnopqrstuvwxyz"),
		Upper: []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
	}

	abcRussian = alphabet{
		Lower: []rune("абвгдеёжзийклмнопрстуфхцчшщъыьэюя"),
		Upper: []rune("АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"),
	}

	abcUkrainian = alphabet{
		Lower: []rune("абвгґдеєжзиіїйклмнопрстуфхцчшщьюя"),
		Upper: []rune("АБВГҐДЕЄЖЗИІЇЙКЛМНОПРСТУФХЦЧШЩЬЮЯ"),
	}
)

func randAbcWord(r *rand.Rand, abc *alphabet) string {

	const maxLetterCount = 12

	rs := make([]rune, r.Intn(maxLetterCount))

	var (
		lower = abc.Lower
		upper = abc.Upper
	)

	for i := range rs {

		if random.Bool(r) {
			rs[i] = lower[r.Intn(len(lower))]
		} else {
			rs[i] = upper[r.Intn(len(upper))]
		}
	}

	return string(rs)
}

func randWord(r *rand.Rand) string {

	var s string

	switch n := r.Intn(3); n {

	case 0:
		s = randAbcWord(r, &abcEnglish)

	case 1:
		s = randAbcWord(r, &abcRussian)

	case 2:
		s = randAbcWord(r, &abcUkrainian)
	}

	return s
}
