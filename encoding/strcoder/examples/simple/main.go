package main

import (
	"fmt"

	"github.com/envoker/golang/encoding/strcoder"
)

func main() {

	unicodeTest()

	//		e := strcoder.NewEncoder("cp1251")
	//		strSrc := "АаБб ... ЮюЯяЇїЄєЁёҐґ  汉字 Юникод (Unicode) — универсальные кодировки UTF 8, 16 и 32"
	//		data := e.Encode(strSrc)
	//		fmt.Printf("% x\n", data)
	//		d := strcoder.NewDecoder("Windows-1251")
	//		strDst := d.Decode(data)
	//		fmt.Println(strDst)
}

func unicodeTest() error {

	const codec = "UTF-16BE"

	var (
		e = strcoder.NewEncoder(codec)
		d = strcoder.NewDecoder(codec)
	)

	text := []byte("Hello, 世界, Привет!")

	p := e.Encode(text)

	fmt.Printf("utf8 enc: % x\n", text)
	fmt.Printf("utf16 enc: % x\n", p)

	q := d.Decode(p)
	fmt.Println(q)

	return nil
}

func tablePrint() {

	bs := make([]byte, 256)

	for i := range bs {
		bs[i] = byte(i)
	}

	//d := strcoder.NewDecoder("cp866")
	d := strcoder.NewAsciiDecoder(strcoder.CP866)

	s := d.Decode(bs)
	fmt.Println(s)
}
