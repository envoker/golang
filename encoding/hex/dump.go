package hex

import (
	"bytes"
)

func Dump(data []byte) string {

	buffer := new(bytes.Buffer)
	for i, b := range data {

		// write separator
		switch {
		case (i == 0):
			// skip
		case (i%128 == 0):
			buffer.WriteString("\n\n")
		case (i%16 == 0):
			buffer.WriteRune('\n')
		case (i%8 == 0):
			buffer.WriteString("  ")
		default:
			buffer.WriteRune(byteSeparator)
		}

		// write byte
		buffer.WriteByte(hexLower[hiNibble(b)])
		buffer.WriteByte(hexLower[loNibble(b)])
	}
	return buffer.String()
}
