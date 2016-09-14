package wav

import (
	"bytes"
	"encoding/binary"
)

var byteOrder = binary.LittleEndian

const (
	size_Uint16      = 2
	size_Uint32      = 4
	size_chunkID     = size_Uint32
	size_chunkSize   = size_Uint32
	size_chunkHeader = size_chunkID + size_chunkSize
	size_Wave        = size_Uint32
	size_RiffHeader  = size_chunkHeader + size_Wave
	size_FmtData     = (4 * size_Uint16) + (2 * size_Uint32)
	size_FmtChunk    = size_chunkHeader + size_FmtData
)

type chunkID [size_chunkID]byte

func (a chunkID) Equal(b chunkID) bool {
	return bytes.Equal(a[:], b[:])
}

var (
	token_RIFF = chunkID{'R', 'I', 'F', 'F'}
	token_WAVE = chunkID{'W', 'A', 'V', 'E'}
	token_fmt  = chunkID{'f', 'm', 't', ' '}
	token_data = chunkID{'d', 'a', 't', 'a'}
)

type chunkHeader struct {
	id   chunkID
	size uint32
}

func (ch *chunkHeader) chunkIdEqual(id chunkID) bool {
	return ch.id.Equal(id)
}

func (chunkHeader) Size() int {
	return size_chunkHeader
}

func (ch *chunkHeader) encode(data []byte) (n int, err error) {

	if len(data) < size_chunkHeader {
		return 0, ErrorWrongDataLen
	}

	copy(data[0:4], ch.id[0:4])
	byteOrder.PutUint32(data[4:8], ch.size)

	return size_chunkHeader, nil
}

func (ch *chunkHeader) decode(data []byte) (n int, err error) {

	if len(data) < size_chunkHeader {
		return 0, ErrorWrongDataLen
	}

	copy(ch.id[0:4], data[0:4])
	ch.size = byteOrder.Uint32(data[4:8])

	return size_chunkHeader, nil
}

type Config struct {
	AudioFormat    int // тип формата (1 - PCM; 6 - A-law, 7 - Mu-law)
	Channels       int // количество каналов (1 - моно; 2 - стeрео)
	SampleRate     int // частота дискретизации (8000, ...)
	BytesPerSample int // 1, 2, 3, 4
}

func (c *Config) checkError() error {

	if (c.AudioFormat < 0) || (c.AudioFormat > 65535) {
		return ErrAudioFormat
	}

	if (c.Channels < 1) || (c.Channels > 32) {
		return ErrChannels
	}

	if (c.SampleRate < 10) || (c.SampleRate > 200000) {
		return ErrSampleRate
	}

	switch c.BytesPerSample {
	case 1, 2, 3, 4:
	default:
		return ErrBytesPerSample
	}

	return nil
}

func (c *Config) BytesPerSec() int {
	return c.Channels * c.BytesPerSample * c.SampleRate
}

func (c *Config) BytesPerBlock() int {
	return c.Channels * c.BytesPerSample
}

type fmtData struct {
	AudioFormat   uint16
	Channels      uint16
	SampleRate    uint32
	BytesPerSec   uint32
	BytesPerBlock uint16
	BitsPerSample uint16
}

func (d *fmtData) setConfig(c Config) error {

	if err := c.checkError(); err != nil {
		return err
	}

	*d = fmtData{
		AudioFormat:   uint16(c.AudioFormat),
		Channels:      uint16(c.Channels),
		SampleRate:    uint32(c.SampleRate),
		BitsPerSample: uint16(c.BytesPerSample * 8),
		BytesPerSec:   uint32(c.BytesPerSec()),
		BytesPerBlock: uint16(c.BytesPerBlock()),
	}

	return nil
}

func (d *fmtData) getConfig(c *Config) {
	c.AudioFormat = int(d.AudioFormat)
	c.Channels = int(d.Channels)
	c.SampleRate = int(d.SampleRate)
	c.BytesPerSample = int(d.BitsPerSample) / 8
}

func (fmtData) Size() int {
	return size_FmtData
}

func (d *fmtData) encode(data []byte) (n int, err error) {

	if len(data) < size_FmtData {
		err = newError("wave config encode: wrong data len")
		return
	}

	byteOrder.PutUint16(data[0:2], d.AudioFormat)
	byteOrder.PutUint16(data[2:4], d.Channels)
	byteOrder.PutUint32(data[4:8], d.SampleRate)
	byteOrder.PutUint32(data[8:12], d.BytesPerSec)
	byteOrder.PutUint16(data[12:14], d.BytesPerBlock)
	byteOrder.PutUint16(data[14:16], d.BitsPerSample)

	return size_FmtData, nil
}

func (d *fmtData) decode(data []byte) (n int, err error) {

	if len(data) < size_FmtData {
		err = newError("wave config decode: wrong data len")
		return
	}

	d.AudioFormat = byteOrder.Uint16(data[0:2])
	d.Channels = byteOrder.Uint16(data[2:4])
	d.SampleRate = byteOrder.Uint32(data[4:8])
	d.BytesPerSec = byteOrder.Uint32(data[8:12])
	d.BytesPerBlock = byteOrder.Uint16(data[12:14])
	d.BitsPerSample = byteOrder.Uint16(data[14:16])

	return size_FmtData, nil
}
