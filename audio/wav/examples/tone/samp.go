package main

import (
	"bufio"
	"time"

	"github.com/envoker/golang/audio/wav"
)

func GenerateWave(fileName string, duration time.Duration, sampleRate float32, bytesPerSample int, samplers []wav.NextSampler) error {

	Tmax := float32(duration.Seconds())

	c := wav.Config{
		AudioFormat:    wav.WAVE_FORMAT_PCM,
		Channels:       uint16(len(samplers)),
		SampleRate:     uint32(sampleRate),
		BytesPerSample: uint16(bytesPerSample),
	}

	fw, err := wav.OpenFileWriter(fileName, &c)
	if err != nil {
		return err
	}
	defer fw.Close()

	bw := bufio.NewWriterSize(fw, int(c.BytesPerSec()))

	w, err := wav.NewSampleWriter(bw, int(c.BytesPerSample))
	if err != nil {
		return err
	}

	n := int(Tmax * sampleRate)
	for i := 0; i < n; i++ {
		for _, sampler := range samplers {
			err = w.WriteSample(sampler.NextSample())
			if err != nil {
				return err
			}
		}
	}

	if err = bw.Flush(); err != nil {
		return err
	}

	return nil
}
