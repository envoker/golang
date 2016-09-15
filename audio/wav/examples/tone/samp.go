package main

import (
	"bufio"
	"time"

	"github.com/envoker/golang/audio/wav"
	"github.com/envoker/golang/audio/wav/sample"
)

func GenerateWave(fileName string, duration time.Duration, sampleRate float32, bytesPerSample int, samplers []sample.NextSampler) error {

	Tmax := float32(duration.Seconds())

	c := wav.Config{
		AudioFormat:    wav.WAVE_FORMAT_PCM,
		Channels:       len(samplers),
		SampleRate:     int(sampleRate),
		BytesPerSample: bytesPerSample,
	}

	fw, err := wav.NewFileWriter(fileName, c)
	if err != nil {
		return err
	}
	defer fw.Close()

	bw := bufio.NewWriterSize(fw, int(c.BytesPerSec()))

	sw, err := sample.NewSampleWriter(bw, int(c.BytesPerSample))
	if err != nil {
		return err
	}

	n := int(Tmax * sampleRate)
	for i := 0; i < n; i++ {
		for _, sampler := range samplers {
			if err = sw.WriteSample(sampler.NextSample()); err != nil {
				return err
			}
		}
	}

	return bw.Flush()
}
