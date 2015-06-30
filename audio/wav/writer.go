package wav

import (
	"io"
	"os"
)

type fileWriter struct {
	config     Config
	dataLength uint32
	f          *os.File
}

func OpenFileWriter(fileName string, config *Config) (io.WriteCloser, error) {

	if config == nil {
		return nil, newError("OpenFileWriter: config is nil")
	}

	if err := config.Error(); err != nil {
		return nil, err
	}

	f, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}

	fw := &fileWriter{
		config:     *config,
		dataLength: 0,
		f:          f,
	}

	fw.writeConfig()

	return fw, nil
}

func (this *fileWriter) Write(data []byte) (n int, err error) {

	n, err = this.f.Write(data)
	if err != nil {
		return
	}

	this.dataLength += uint32(n)

	return
}

func (this *fileWriter) Close() error {

	if this.f != nil {
		this.writeDataLength()

		this.f.Close()
		this.f = nil
		this.dataLength = 0
	}

	return nil
}

func (this *fileWriter) writeConfig() error {

	_, err := this.f.Seek(0, os.SEEK_SET)
	if err != nil {
		return err
	}

	var ch chunkHeader

	// RIFF header
	{
		ch = chunkHeader{
			id:   token_RIFF,
			size: 0,
		}

		if _, err = encodeAndWrite(&ch, this.f); err != nil {
			return err
		}

		// WAVE
		{
			var riffFormat = token_WAVE
			if _, err = this.f.Write(riffFormat[:]); err != nil {
				return err
			}
		}
	}

	// fmt chunk
	{
		ch = chunkHeader{
			id:   token_fmt,
			size: size_FmtData,
		}

		if _, err = encodeAndWrite(&ch, this.f); err != nil {
			return err
		}

		var c_data fmtData
		if err = c_data.setConfig(this.config); err != nil {
			return err
		}

		if _, err = encodeAndWrite(&c_data, this.f); err != nil {
			return err
		}
	}

	// data chunk header
	{
		ch = chunkHeader{
			id:   token_data,
			size: 0,
		}

		if _, err = encodeAndWrite(&ch, this.f); err != nil {
			return err
		}
	}

	return nil
}

func (this *fileWriter) writeDataLength() (err error) {

	var data = make([]byte, size_chunkSize)

	size_dataChunk := size_chunkHeader + this.dataLength

	// RIFF chunk
	{
		size := uint32(size_Wave + size_FmtChunk + size_dataChunk)
		pos := int64(size_chunkID)

		byteOrder.PutUint32(data, size)
		this.f.Seek(pos, os.SEEK_SET)
		if _, err = this.f.Write(data); err != nil {
			return
		}
	}

	// data chunk
	{
		size := this.dataLength
		pos := int64(size_RiffHeader + size_FmtChunk + size_chunkID)

		byteOrder.PutUint32(data, size)
		this.f.Seek(pos, os.SEEK_SET)
		if _, err = this.f.Write(data); err != nil {
			return
		}
	}

	return
}
