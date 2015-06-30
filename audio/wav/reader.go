package wav

import (
	"io"
	"os"
)

type fileReader struct {
	config     Config
	dataLength uint32
	f          *os.File
}

func OpenFileReader(fileName string, config *Config) (io.ReadCloser, error) {

	if config == nil {
		return nil, newError("OpenFileReader: config is nil")
	}

	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	fr := &fileReader{
		dataLength: 0,
		f:          f,
	}

	if err = fr.readConfig(); err != nil {
		return nil, err
	}

	*config = fr.config

	return fr, nil
}

func (this *fileReader) Read(data []byte) (n int, err error) {

	if this.dataLength == 0 {
		return 0, nil
	}

	n = len(data)
	if n > int(this.dataLength) {
		n = int(this.dataLength)
	}

	n, err = this.f.Read(data[:n])
	if err != nil {
		return 0, err
	}

	this.dataLength -= uint32(n)

	return
}

func (this *fileReader) Close() error {

	if this.f != nil {
		this.f.Close()
		this.f = nil
		this.dataLength = 0
	}

	return nil
}

func (this *fileReader) getConfig(c *Config) error {

	*c = this.config

	return nil
}

func (this *fileReader) readConfig() error {

	var err error

	_, err = this.f.Seek(0, os.SEEK_SET)
	if err != nil {
		return err
	}

	var riffSize uint32
	var ch chunkHeader

	// RIFF header
	{
		if _, err = readAndDecode(this.f, &ch); err != nil {
			return err
		}

		if !ch.chunkIdEqual(token_RIFF) {
			return ErrorItIsNotRiffFile
		}

		riffSize = ch.size

		// WAVE
		{
			var riffFormat chunkID

			if _, err = this.f.Read(riffFormat[:]); err != nil {
				return err
			}

			if !riffFormat.Equal(token_WAVE) {
				return ErrorItIsNotRiffFile
			}
		}
	}

	var f_fmtChunk, f_dataChunk bool
	var n_riffSize = uint32(size_Wave)

	ever := true
	for ever {

		if _, err = readAndDecode(this.f, &ch); err != nil {
			return err
		}

		n_riffSize += size_chunkHeader + ch.size

		switch {

		case ch.chunkIdEqual(token_fmt):
			{
				var c_data fmtData

				_, err = readAndDecode(this.f, &c_data)
				if err != nil {
					return err
				}

				this.config, err = c_data.getConfig()
				if err != nil {
					return err
				}

				f_fmtChunk = true
			}

		case ch.chunkIdEqual(token_data):
			{
				this.dataLength = ch.size
				f_dataChunk = true
				ever = false
				break
			}

		default: // skip other chunk data
			{
				_, err = this.f.Seek(int64(ch.size), os.SEEK_CUR)
				if err != nil {
					return err
				}
			}
		}
	}

	if (!f_fmtChunk) || (!f_dataChunk) {
		err = ErrorChunkStructure
		return err
	}

	if n_riffSize != riffSize {
		err = ErrorChunkStructure
		return err
	}

	return nil
}
