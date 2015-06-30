package wav

import (
	"io"
	"math"
)

const (
	maxInt8 = (1 << ((iota+1)*8 - 1)) - 1
	maxInt16
	maxInt24
	maxInt32
	maxInt40
	maxInt48
	maxInt56
	maxInt64
)

func floor(x float32) float32 {
	return float32(math.Floor(float64(x)))
}

func round(x float32) float32 {
	return floor(x + 0.5)
}

func SampleNormalize(sample float32) float32 {

	if sample > 1.0 {
		sample = 1.0
	}
	if sample < -1.0 {
		sample = -1.0
	}

	return sample
}

/*
type Sample float32

func (s *Sample) normalize() {

	switch {

	case (*s > 1.0):
		*s = 1.0

	case (*s < -1.0):
		*s = -1.0
	}
}

func (s Sample) Int8() int8 {
	s.normalize()
	return int8(round(maxInt8 * float32(s)))
}

func (s Sample) Int16() int16 {
	s.normalize()
	return int16(round(maxInt16 * float32(s)))
}

func (s Sample) Int24() int32 {
	s.normalize()
	return int32(round(maxInt24 * float32(s)))
}

func (s Sample) Int32() int32 {
	s.normalize()
	return int32(round(maxInt32 * float32(s)))
}

func (s Sample) Int48() int64 {
	s.normalize()
	return int64(round(maxInt48 * float32(s)))
}

func (s Sample) Int64() int64 {
	s.normalize()
	return int64(round(maxInt64 * float32(s)))
}
*/

//-----------------------------------------
func SampleToInt8(sample float32) int8 {
	sample = SampleNormalize(sample)
	return int8(round(maxInt8 * sample))
}

func SampleToInt16(sample float32) int16 {
	sample = SampleNormalize(sample)
	return int16(round(maxInt16 * sample))
}

func SampleToInt24(sample float32) int32 {
	sample = SampleNormalize(sample)
	return int32(round(maxInt24 * sample))
}

func SampleToInt32(sample float32) int32 {
	sample = SampleNormalize(sample)
	return int32(round(maxInt32 * sample))
}

func SampleToInt48(sample float32) int64 {
	sample = SampleNormalize(sample)
	return int64(round(maxInt48 * sample))
}

func SampleToInt64(sample float32) int64 {
	sample = SampleNormalize(sample)
	return int64(round(maxInt64 * sample))
}

//-----------------------------------------
func SampleFromInt8(i int8) float32 {
	return SampleNormalize(float32(i) / maxInt8)
}

func SampleFromInt16(i int16) float32 {
	return SampleNormalize(float32(i) / maxInt16)
}

func SampleFromInt24(i int32) float32 {
	return SampleNormalize(float32(i) / maxInt24)
}

func SampleFromInt32(i int32) float32 {
	return SampleNormalize(float32(i) / maxInt32)
}

func SampleFromInt48(i int64) float32 {
	return SampleNormalize(float32(i) / maxInt48)
}

func SampleFromInt64(i int64) float32 {
	return SampleNormalize(float32(i) / maxInt64)
}

//-----------------------------------------
type sampler interface {
	writeSample(w io.Writer, sample float32) error
}

type int8Sampler struct {
	bs [1]byte
}

func (this *int8Sampler) writeSample(w io.Writer, sample float32) error {

	b := this.bs[:]

	i := SampleToInt8(sample)
	u := uint8(int(i) + 128)
	b[0] = byte(u)

	if _, err := w.Write(b); err != nil {
		return err
	}

	return nil
}

func (this *int8Sampler) readSample(r io.Reader) (sample float32, err error) {

	b := this.bs[:]

	if _, err = r.Read(b); err != nil {
		return
	}

	u := uint8(b[0])
	i := int8(int(u) - 128)
	sample = SampleFromInt8(i)

	return
}

type int16Sampler struct {
	bs [2]byte
}

func (this *int16Sampler) writeSample(w io.Writer, sample float32) error {

	b := this.bs[:]

	i := SampleToInt16(sample)
	le.PutInt16(b, i)

	if _, err := w.Write(b); err != nil {
		return err
	}

	return nil
}

func (this *int16Sampler) readSample(r io.Reader) (sample float32, err error) {

	b := this.bs[:]

	if _, err = r.Read(b); err != nil {
		return
	}

	i := le.GetInt16(b)
	sample = SampleFromInt16(i)

	return
}

type int24Sampler struct {
	bs [3]byte
}

func (this *int24Sampler) writeSample(w io.Writer, sample float32) error {

	b := this.bs[:]

	i := SampleToInt24(sample)
	le.PutInt24(b, i)

	if _, err := w.Write(b); err != nil {
		return err
	}

	return nil
}

type int32Sampler struct {
	bs [4]byte
}

func (this *int32Sampler) writeSample(w io.Writer, sample float32) error {

	b := this.bs[:]

	i := SampleToInt32(sample)
	le.PutInt32(b, i)

	if _, err := w.Write(b); err != nil {
		return err
	}

	return nil
}

type NextSampler interface {
	NextSample() float32 // [ -1 ... +1 ]
}

/*
type SampleWriter interface {
	WriteSample(sample float32) (err error)
}

type SampleReader interface {
	ReadSample() (sample float32, err error)
}
*/

type SampleWriter struct {
	w io.Writer
	s sampler
}

func (this *SampleWriter) WriteSample(sample float32) error {

	return this.s.writeSample(this.w, sample)
}

func NewSampleWriter(w io.Writer, bytesPerSample int) (sw *SampleWriter, err error) {

	switch bytesPerSample {

	case 1:
		sw = &SampleWriter{w, new(int8Sampler)}

	case 2:
		sw = &SampleWriter{w, new(int16Sampler)}

	case 3:
		sw = &SampleWriter{w, new(int24Sampler)}

	case 4:
		sw = &SampleWriter{w, new(int32Sampler)}

	default:
		err = ErrorBytesPerSample
	}

	return
}

//---------------------------------------------------------------------------------
/*
type BlockWriter interface {
	WriteBlock(block []float32) error
}

type BlockReader interface {
	ReadBlock(block []float32) error
}

//---------------------------------------------------------------------------------
type blockWriter struct {
	w  io.Writer
	c  Config
	sw SampleWriter
}

func NewBlockWriter(w io.Writer, c Config) (BlockWriter, error) {

	if err := c.Error(); err != nil {
		return nil, err
	}

	sw, err := NewSampleWriter(w, int(c.BytesPerSample))
	if err != nil {
		return nil, err
	}

	bw := &blockWriter{
		w:  w,
		c:  c,
		sw: sw,
	}

	return bw, nil
}

func (this *blockWriter) WriteBlock(block []float32) (err error) {

	if int(this.c.Channels) != len(block) {
		err = newError("wrong Channels")
		return
	}

	for _, sample := range block {
		if err = this.sw.WriteSample(sample); err != nil {
			return
		}
	}

	return
}
*/