package alaw

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
)

func ConvertFile(sourceName, destName string) {

	var (
		err          error
		n            int
		linearSample int16
	)

	var sourceFile, destFile *os.File

	sourceFile, err = os.Open(sourceName)
	if err != nil {
		fmt.Println("error: open source file: ", sourceName)
		return
	}
	defer sourceFile.Close()
	r := bufio.NewReader(sourceFile)

	destFile, err = os.Create(destName)
	if err != nil {
		fmt.Println("error: open dest file: ", destName)
		return
	}
	defer destFile.Close()
	w := bufio.NewWriter(destFile)

	const blen = 1024
	buff := make([]byte, blen)
	buffLinear := make([]byte, 2*blen)

	for {
		n, err = r.Read(buff)

		if err != nil {
			break
		}

		for i := 0; i < n; i++ {
			linearSample = ALawToLinear(buff[i])

			j := i * 2
			binary.BigEndian.PutUint16(buffLinear[j:], uint16(linearSample))
		}

		w.Write(buffLinear[:2*n])
	}

	w.Flush()
}
