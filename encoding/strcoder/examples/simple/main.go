package main

import (
	"fmt"

	"github.com/envoker/golang/encoding/strcoder"
)

func main() {

	tablePrint()

	e := strcoder.NewEncoder("cp1251")

	strSrc := "АаБб ... ЮюЯяЇїЄєЁёҐґ  汉字 Юникод (Unicode) — универсальные кодировки UTF 8, 16 и 32"

	data := e.Encode(strSrc)
	fmt.Printf("% x\n", data)

	d := strcoder.NewDecoder("Windows-1251")

	strDst := d.Decode(data)
	fmt.Println(strDst)
}

func tablePrint() {

	bs := make([]byte, 256)

	for i := range bs {
		bs[i] = byte(i)
	}

	//d := strcoder.NewDecoder("cp866")
	d := strcoder.NewAsciiDecoder(strcoder.Cp866)

	s := d.Decode(bs)
	fmt.Println(s)
}
