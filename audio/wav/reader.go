package wav

import "os"

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

	var riffSize uint32
	var ch chunkHeader

	// RIFF header
	{
		if _, err = readAndDecode(fr.file, &ch); err != nil {
			return err
		}

		if !ch.chunkIdEqual(token_RIFF) {
			return ErrorItIsNotRiffFile
		}

		riffSize = ch.size

		// WAVE
		{
			var riffFormat chunkID

			if _, err = fr.file.Read(riffFormat[:]); err != nil {
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

		if _, err = readAndDecode(fr.file, &ch); err != nil {
			return err
		}

		n_riffSize += size_chunkHeader + ch.size

		switch {

		case ch.chunkIdEqual(token_fmt):
			{
				var c_data fmtData

				_, err = readAndDecode(fr.file, &c_data)
				if err != nil {
					return err
				}

				c_data.getConfig(&(fr.config))

				f_fmtChunk = true
			}

		case ch.chunkIdEqual(token_data):
			{
				fr.dataLength = ch.size
				f_dataChunk = true
				ever = false
				break
			}

		default: // skip other chunk data
			{
				_, err = fr.file.Seek(int64(ch.size), os.SEEK_CUR)
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
