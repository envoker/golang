package json

import (
	"unicode/utf8"
)

const (
	rc_FormFeed       = '\f'
	rc_CarriageReturn = '\r'
	rc_NewLine        = '\n'
	rc_HorizontalTab  = '\t'
	rc_VerticalTab    = '\v'
	rc_Backslash      = '\\'
	rc_SingleQuote    = '\''

	rc_DoubleQuotes = '"'
	rc_Slash        = '/'
	rc_Space        = ' '
	rc_Plus         = '+'
	rc_Minus        = '-'
	rc_Comma        = ','
	rc_Dot          = '.'
	rc_Colon        = ':'

	rc_OpenCurlyBracket   = '{'
	rc_CloseCurlyBracket  = '}'
	rc_OpenSquareBracket  = '['
	rc_CloseSquareBracket = ']'

	rc_Backspace = '\x08'

	rc_NumberZero  = '0'
	rc_NumberOne   = '1'
	rc_NumberTwo   = '2'
	rc_NumberThree = '3'
	rc_NumberFour  = '4'
	rc_NumberFive  = '5'
	rc_NumberSix   = '6'
	rc_NumberSeven = '7'
	rc_NumberEight = '8'
	rc_NumberNine  = '9'
)

const (
	mask_IsSpace uint16 = 1 << iota
	mask_IsDigit
	mask_IsNullBegin
	mask_IsBooleanBegin
	mask_IsNumberBegin
	mask_IsNull
	mask_IsBoolean
	mask_IsNumber
)

var (
	vs_bool = []string{
		"true", "TRUE", "True",
		"false", "FALSE", "False",
	}
	vs_null = []string{"null", "NULL", "Null"}
)

type checkTable [256]uint16

func newCheckTable() *checkTable {

	table := new(checkTable)
	for i := range table {

		r := rune(i)
		var mask uint16

		if slow_IsSpace(r) {
			mask |= mask_IsSpace
		}
		if slow_IsDigit(r) {
			mask |= mask_IsDigit
		}
		if slow_IsNullBegin(r) {
			mask |= mask_IsNullBegin
		}
		if slow_IsBooleanBegin(r) {
			mask |= mask_IsBooleanBegin
		}
		if slow_IsNumberBegin(r) {
			mask |= mask_IsNumberBegin
		}
		if slow_IsNull(r) {
			mask |= mask_IsNull
		}
		if slow_IsBoolean(r) {
			mask |= mask_IsBoolean
		}
		if slow_IsNumber(r) {
			mask |= mask_IsNumber
		}

		table[i] = mask
	}

	return table
}

func (this *checkTable) check(r rune, mask uint16) (ok bool) {

	if (r >= 0) && (r < 256) {
		ok = ((this[r] & mask) != 0)
	}

	return
}

func (this *checkTable) IsNullBegin(r rune) bool {

	return this.check(r, mask_IsNullBegin)
}

func (this *checkTable) IsBooleanBegin(r rune) bool {

	return this.check(r, mask_IsBooleanBegin)
}

func (this *checkTable) IsNumberBegin(r rune) bool {

	return this.check(r, mask_IsNumberBegin)
}

func (this *checkTable) IsSpace(r rune) bool {

	return this.check(r, mask_IsSpace)
}

func (this *checkTable) IsDigit(r rune) bool {

	return this.check(r, mask_IsDigit)
}

func (this *checkTable) IsNull(r rune) bool {

	return this.check(r, mask_IsNull)
}

func (this *checkTable) IsBoolean(r rune) bool {

	return this.check(r, mask_IsBoolean)
}

func (this *checkTable) IsNumber(r rune) bool {

	return this.check(r, mask_IsNumber)
}

//---------------------------------------------------------------------------------
var ct = newCheckTable()

func runeIn(rs []rune, r rune) bool {

	for _, p := range rs {
		if p == r {
			return true
		}
	}

	return false
}

//---------------------------------------------------------------------------------
//	Slow functions
//---------------------------------------------------------------------------------
func slow_IsSpace(r rune) bool {

	var rs = []rune{
		rc_HorizontalTab,
		rc_NewLine,
		rc_VerticalTab,
		rc_FormFeed,
		rc_CarriageReturn,
		rc_Space,
	}

	return runeIn(rs, r)
}

func slow_IsDigitSign(r rune) bool {

	var rs = []rune{rc_Plus, rc_Minus}
	return runeIn(rs, r)
}

func slow_IsDigit(r rune) bool {

	var rs = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	return runeIn(rs, r)
}

func slow_IsNullBegin(r rune) (ok bool) {

	return runeIn(firstRunes(vs_null), r)
}

func slow_IsBooleanBegin(r rune) (ok bool) {

	return runeIn(firstRunes(vs_bool), r)
}

func slow_IsNumberBegin(r rune) bool {

	if slow_IsDigit(r) {
		return true
	}

	if slow_IsDigitSign(r) {
		return true
	}

	return false
}

func slow_IsNull(r rune) (ok bool) {

	return runeIn(uniqueRunes(vs_null), r)
}

func slow_IsBoolean(r rune) (ok bool) {

	return runeIn(uniqueRunes(vs_bool), r)
}

func slow_IsNumber(r rune) bool {

	if slow_IsDigit(r) {
		return true
	}

	if slow_IsDigitSign(r) {
		return true
	}

	var rs = []rune{rc_Dot, 'e', 'E'}
	return runeIn(rs, r)
}

func uniqueRunes(strs []string) []rune {

	rs := make([]rune, 0)

	for _, s := range strs {

		for {
			r, size := utf8.DecodeRuneInString(s)
			if size == 0 {
				break
			}
			s = s[size:]

			f_append := true
			for i := range rs {
				if rs[i] == r {
					f_append = false
					break
				}
			}

			if f_append {
				rs = append(rs, r)
			}
		}
	}

	return rs
}

func firstRunes(strs []string) []rune {

	rs := make([]rune, 0)

	for _, s := range strs {

		r, size := utf8.DecodeRuneInString(s)
		if size == 0 {
			continue
		}

		f_append := true
		for i := range rs {
			if rs[i] == r {
				f_append = false
				break
			}
		}

		if f_append {
			rs = append(rs, r)
		}
	}

	return rs
}
