package main

import (
	"fmt"
	"os"

	"github.com/envoker/golang/audio/wav"
	chex "github.com/envoker/golang/encoding/hex"
)

func main() {

	wav.ErrorBytesPerSample

	fileName := "./test.wav"

	TestWaveRead(fileName)
	//TestWaveWrite(fileName)
	TextFileHexDump(fileName)
}

func TestWaveWrite(fileName string) {

	c := wav.Config{
		AudioFormat:    wav.WAVE_FORMAT_PCM,
		Channels:       2,
		SampleRate:     22050,
		BytesPerSample: 2,
	}

	fw, err := wav.NewFileWriter(fileName, c)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fw.Close()

	d := make([]byte, 2048)

	for i := range d {
		d[i] = byte(i)
	}

	fw.Write(d)
}

func TestWaveRead(fileName string) {

	var c wav.Config

	fr, err := wav.NewFileReader(fileName, &c)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fr.Close()

	fmt.Printf("%+v\n", c)
}

func TextFileHexDump(fileName string) error {

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	data := make([]byte, 100)

	n, err := f.Read(data)
	if err != nil {
		return err
	}

	h := chex.Dump(data[:n])
	fmt.Println(h)

	return nil
}
