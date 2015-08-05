package main

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"log"
	"os"

	"github.com/envoker/golang/encoding/strcoder"
)

func main() {

	//-dest="./text_utf8.txt" -src="./text_cp1251.txt" -codec="cp1251"

	pDest := flag.String("dest", "", "destination file name")
	pSrc := flag.String("src", "", "source file name")
	pCodec := flag.String("codec", "cp1251", "codec type")

	flag.Parse()

	if *pDest == "" {
		log.Fatal("Absent dest")
	}
	if *pSrc == "" {
		log.Fatal("Absent src")
	}

	if err := fileDecode(*pDest, *pSrc, *pCodec); err != nil {
		log.Fatal(err)
	}
}

func fileDecode(dest, src string, codec string) error {

	d := strcoder.NewDecoder(codec)
	if d == nil {
		return errors.New("wrong codec")
	}

	fr, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fr.Close()
	r := bufio.NewReader(fr)

	fw, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer fw.Close()
	w := bufio.NewWriter(fw)
	defer w.Flush()

	data := make([]byte, 1024)
	for {
		n, err := r.Read(data)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		s := d.Decode(data[:n])

		if _, err = w.WriteString(s); err != nil {
			return err
		}
	}

	return nil
}
