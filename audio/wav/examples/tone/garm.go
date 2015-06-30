package main

import (
	"math"

	"github.com/envoker/golang/audio/wav"
)

type Garmonica struct {
	Amplitude, Frequency, Phase float32
}

type toneSampler struct {
	amplitude float32
	phase     float32
	w         float32
	t, dt     float32
}

func (this *toneSampler) NextSample() float32 {

	u := float64(this.w*this.t + this.phase)
	sample := this.amplitude * float32(math.Sin(u))

	this.t += this.dt

	return sample
}

func NewToneSampler(g Garmonica, sampleRate float32) wav.NextSampler {

	return &toneSampler{
		amplitude: g.Amplitude,
		phase:     g.Phase,
		w:         2 * math.Pi * g.Frequency,
		t:         0,
		dt:        1.0 / sampleRate,
	}
}

func MakeSamplers(gs []Garmonica, sampleRate float32) []wav.NextSampler {

	samplers := make([]wav.NextSampler, len(gs))
	for i, g := range gs {
		samplers[i] = NewToneSampler(g, sampleRate)
	}

	return samplers
}
