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

func (this chunkID) Equal(other chunkID) bool {
	return bytes.Equal(this[:], other[:])
}

var (
	token_RIFF = chunkID{'R', 'I', 'F', 'F'}
	token_WAVE = chunkID{'W', 'A', 'V', 'E'}
	token_fmt  = chunkID{'f', 'm', 't', ' '}
	token_data = chunkID{'d', 'a', 't', 'a'}
)

//---------------------------------------------------------
type chunkHeader struct {
	id   chunkID
	size uint32
}

func (this *chunkHeader) chunkIdEqual(id chunkID) bool {
	return this.id.Equal(id)
}

func (this *chunkHeader) Size() int {
	return size_chunkHeader
}

func (this *chunkHeader) encode(data []byte) (n int, err error) {

	if len(data) < size_chunkHeader {
		err = ErrorWrongDataLen
		return
	}

	copy(data[0:4], this.id[0:4])
	byteOrder.PutUint32(data[4:8], this.size)

	return size_chunkHeader, nil
}

func (this *chunkHeader) decode(data []byte) (n int, err error) {

	if len(data) < size_chunkHeader {
		err = ErrorWrongDataLen
		return
	}

	copy(this.id[0:4], data[0:4])
	this.size = byteOrder.Uint32(data[4:8])

	return size_chunkHeader, nil
}

//--------------------------------------------------------------------------
type Config struct {
	AudioFormat    uint16 // тип формата (1 - PCM; 6 - A-law, 7 - Mu-law)
	Channels       uint16 // количество каналов (1 - моно; 2 - стeрео)
	SampleRate     uint32 // частота дискретизации (8000, ...)
	BytesPerSample uint16 // 1, 2, 3, 4
}

func (this *Config) Error() error {

	if (this.Channels < 1) || (this.Channels > 32) {
		return newError("Config.Error(): wrong Channels")
	}

	if (this.SampleRate < 10) || (this.SampleRate > 200000) {
		return newError("Config.Error(): wrong SampleRate")
	}

	switch this.BytesPerSample {
	case 1, 2, 3, 4:
	default:
		return newError("Config.Error(): wrong BytesPerSample")
	}

	return nil
}

func (this *Config) IsValid() bool {

	if err := this.Error(); err != nil {
		return false
	}

	return true
}

func (this *Config) BytesPerSec() uint32 {
	return uint32(this.Channels) * uint32(this.BytesPerSample) * uint32(this.SampleRate)
}

func (this *Config) BytesPerBlock() uint16 {
	return uint16(this.Channels) * uint16(this.BytesPerSample)
}

//-------------------------------------------------------
type fmtData struct {
	AudioFormat   uint16
	Channels      uint16
	SampleRate    uint32
	BytesPerSec   uint32
	BytesPerBlock uint16
	BitsPerSample uint16
}

func (this *fmtData) setConfig(c Config) (err error) {

	if !c.IsValid() {
		err = newError("config is not valid")
		return
	}

	this.AudioFormat = c.AudioFormat
	this.Channels = c.Channels
	this.SampleRate = c.SampleRate
	this.BitsPerSample = c.BytesPerSample * 8
	this.BytesPerSec = c.BytesPerSec()
	this.BytesPerBlock = c.BytesPerBlock()

	return
}

func (this *fmtData) getConfig() (c Config, err error) {

	c.AudioFormat = this.AudioFormat
	c.Channels = this.Channels
	c.SampleRate = this.SampleRate
	c.BytesPerSample = this.BitsPerSample / 8

	return
}

func (this *fmtData) Size() int {
	return size_FmtData
}

func (this *fmtData) encode(data []byte) (n int, err error) {

	if len(data) < size_FmtData {
		err = newError("wave config encode: wrong data len")
		return
	}

	byteOrder.PutUint16(data[0:2], this.AudioFormat)
	byteOrder.PutUint16(data[2:4], this.Channels)
	byteOrder.PutUint32(data[4:8], this.SampleRate)
	byteOrder.PutUint32(data[8:12], this.BytesPerSec)
	byteOrder.PutUint16(data[12:14], this.BytesPerBlock)
	byteOrder.PutUint16(data[14:16], this.BitsPerSample)

	return size_FmtData, nil
}

func (this *fmtData) decode(data []byte) (n int, err error) {

	if len(data) < size_FmtData {
		err = newError("wave config decode: wrong data len")
		return
	}

	this.AudioFormat = byteOrder.Uint16(data[0:2])
	this.Channels = byteOrder.Uint16(data[2:4])
	this.SampleRate = byteOrder.Uint32(data[4:8])
	this.BytesPerSec = byteOrder.Uint32(data[8:12])
	this.BytesPerBlock = byteOrder.Uint16(data[12:14])
	this.BitsPerSample = byteOrder.Uint16(data[14:16])

	return size_FmtData, nil
}
