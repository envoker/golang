package json

import (
	"bytes"
	"io"
)

func bw_WriteIndent(bw BufferWriter, indent int) (n int, err error) {

	if indent < 0 {
		err = newError("InvalidIndent")
		return
	}

	var writeCount int
	n = 0
	for i := 0; i < indent; i++ {
		writeCount, err = bw.WriteRune(rc_HorizontalTab)
		if err != nil {
			return
		}
		n += writeCount
	}

	return
}

func bw_WriteEndOfLine(bw BufferWriter) (err error) {

	if _, err = bw.WriteRune(rc_NewLine); err != nil {
		return
	}

	return
}

func br_ReadString(br BufferReader, runeIsValid func(r rune) bool) (string, error) {

	var (
		size int
		r    rune
		err  error
	)

	destBuffer := new(bytes.Buffer)

	for {
		if r, size, err = br.ReadRune(); err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		if size == 0 {
			break
		}

		if !runeIsValid(r) {
			if err = br.UnreadRune(); err != nil {
				return "", err
			}
			break
		}

		if _, err = destBuffer.WriteRune(r); err != nil {
			return "", err
		}
	}

	return string(destBuffer.Bytes()), nil
}

func br_SkipSpaces(br BufferReader) (n int, err error) {

	var (
		r    rune
		size int
	)

	for {

		if r, size, err = br.ReadRune(); err != nil {
			return
		}

		if size == 0 {
			break
		}

		if !ct.IsSpace(r) {
			if err = br.UnreadRune(); err != nil {
				return
			}
			break
		}

		n += size
	}

	return
}

func br_SkipRune(br BufferReader, r rune) (ok bool) {

	var (
		p    rune
		size int
		err  error
	)

	if p, size, err = br.ReadRune(); err != nil {
		return
	}

	if size > 0 {
		if r == p {
			ok = true
		} else {
			if err = br.UnreadRune(); err != nil {
				return
			}
		}
	}

	return
}
