package wav

import "encoding/binary"

var (
	size_chunkId     = 4 // uint32
	size_chunkHeader = binary.Size(chunkHeader{})
	size_FmtData     = binary.Size(fmtData{})

	size_Format     = size_chunkId
	size_RiffHeader = size_chunkHeader + size_Format
	size_FmtChunk   = size_chunkHeader + size_FmtData
)

var (
	token_RIFF uint32 = getToken("RIFF")
	token_WAVE uint32 = getToken("WAVE")
	token_fmt  uint32 = getToken("fmt ")
	token_data uint32 = getToken("data")
)

func getToken(s string) uint32 {
	return binary.LittleEndian.Uint32([]byte(s))
}

type chunkHeader struct {
	Id   uint32
	Size uint32
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

func configToFmtData(c Config) fmtData {
	return fmtData{
		AudioFormat:   uint16(c.AudioFormat),
		Channels:      uint16(c.Channels),
		SampleRate:    uint32(c.SampleRate),
		BitsPerSample: uint16(c.BytesPerSample * 8),
		BytesPerSec:   uint32(c.BytesPerSec()),
		BytesPerBlock: uint16(c.BytesPerBlock()),
	}
}

func fmtDataToConfig(d fmtData) Config {
	return Config{
		AudioFormat:    int(d.AudioFormat),
		Channels:       int(d.Channels),
		SampleRate:     int(d.SampleRate),
		BytesPerSample: int(d.BitsPerSample) / 8,
	}
}
