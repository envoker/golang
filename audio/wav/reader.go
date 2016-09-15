package wav

import (
	"encoding/binary"
	"os"
)

type FileReader struct {
	config     Config
	dataLength uint32
	file       *os.File
}

func NewFileReader(fileName string, config *Config) (*FileReader, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	fr := &FileReader{
		dataLength: 0,
		file:       file,
	}

	if err = fr.readConfig(); err != nil {
		return nil, err
	}

	*config = fr.config

	return fr, nil
}

func (fr *FileReader) Close() error {

	if fr.file == nil {
		return ErrFileReaderClosed
	}

	err := fr.file.Close()
	fr.file = nil
	fr.dataLength = 0
	return err
}

func (fr *FileReader) Read(data []byte) (n int, err error) {

	if fr.dataLength == 0 {
		return 0, nil
	}

	n = len(data)
	if n > int(fr.dataLength) {
		n = int(fr.dataLength)
	}

	n, err = fr.file.Read(data[:n])
	if err != nil {
		return 0, err
	}

	fr.dataLength -= uint32(n)

	return
}

func (fr *FileReader) getConfig(c *Config) {
	*c = fr.config
}

func (fr *FileReader) readConfig() error {

	_, err := fr.file.Seek(0, os.SEEK_SET)
	if err != nil {
		return err
	}

	var (
		riffSize uint32
		ch       chunkHeader
	)

	// RIFF header
	{
		err = binary.Read(fr.file, binary.LittleEndian, &ch)
		if err != nil {
			return err
		}

		if ch.Id != token_RIFF {
			return ErrFileFormat
		}

		riffSize = ch.Size

		// WAVE

		var format uint32

		err = binary.Read(fr.file, binary.LittleEndian, &format)
		if err != nil {
			return err
		}

		if format != token_WAVE {
			return ErrFileFormat
		}
	}

	var f_fmtChunk, f_dataChunk bool
	var n_riffSize = uint32(size_Format)

	ever := true
	for ever {

		err = binary.Read(fr.file, binary.LittleEndian, &ch)
		if err != nil {
			return err
		}

		n_riffSize += uint32(size_chunkHeader) + ch.Size

		switch ch.Id {

		case token_fmt:
			{
				var c_data fmtData

				err = binary.Read(fr.file, binary.LittleEndian, &c_data)
				if err != nil {
					return err
				}

				fr.config = fmtDataToConfig(c_data)

				f_fmtChunk = true
			}

		case token_data:
			{
				fr.dataLength = ch.Size
				f_dataChunk = true
				ever = false
				break
			}

		default: // skip other chunk data
			_, err = fr.file.Seek(int64(ch.Size), os.SEEK_CUR)
			if err != nil {
				return err
			}
		}
	}

	if (!f_fmtChunk) || (!f_dataChunk) {
		return ErrFileFormat
	}

	if n_riffSize != riffSize {
		return ErrFileFormat
	}

	return nil
}
