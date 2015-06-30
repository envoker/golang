package main

import (
	"fmt"
	"math"
	"time"
)

const bytesPerSample = 3

func monoTone() error {

	sampleRate := float32(22050)

	gs := []Garmonica{
		Garmonica{
			Amplitude: 0.9,
			Frequency: 1000,
			Phase:     math.Pi * 0.5,
		},
	}

	samplers := MakeSamplers(gs, sampleRate)

	err := GenerateWave("./mono-tone.wav", time.Second*9, sampleRate, bytesPerSample, samplers)
	if err != nil {
		return err
	}

	return nil
}

func stereoTone() error {

	sampleRate := float32(44100)

	gs := []Garmonica{
		Garmonica{
			Amplitude: 0.75,
			Frequency: 500,
			Phase:     math.Pi * 0.5,
		},
		Garmonica{
			Amplitude: 0.60,
			Frequency: 1000,
			Phase:     0,
		},
	}

	samplers := MakeSamplers(gs, sampleRate)

	err := GenerateWave("./stereo-tone.wav", time.Second*13, sampleRate, bytesPerSample, samplers)
	if err != nil {
		return err
	}

	return nil
}

func multiTone() error {

	sampleRate := float32(44100)

	gs := []Garmonica{
		Garmonica{
			Amplitude: 0.75,
			Frequency: 1500,
			Phase:     math.Pi * 0.5,
		},
		Garmonica{
			Amplitude: 0.60,
			Frequency: 3000,
			Phase:     0,
		},
		Garmonica{
			Amplitude: 0.90,
			Frequency: 400,
			Phase:     0,
		},
		Garmonica{
			Amplitude: 0.50,
			Frequency: 1000,
			Phase:     0,
		},
	}

	samplers := MakeSamplers(gs, sampleRate)

	err := GenerateWave("./multi-tone.wav", time.Second*7, sampleRate, bytesPerSample, samplers)
	if err != nil {
		return err
	}

	return nil
}

func Generate() (err error) {

	if err = monoTone(); err != nil {
		return
	}
	if err = stereoTone(); err != nil {
		return
	}
	if err = multiTone(); err != nil {
		return
	}

	return
}

func main() {

	if err := Generate(); err != nil {
		fmt.Println(err.Error())
	}
}