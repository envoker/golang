package wav

import "os"

type FileWriter struct {
	config     Config
	dataLength uint32
	file       *os.File
}

func NewFileWriter(fileName string, config Config) (*FileWriter, error) {

	if err := config.checkError(); err != nil {
		return nil, err
	}

	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}

	fw := &FileWriter{
		config:     config,
		dataLength: 0,
		file:       file,
	}

	fw.writeConfig()

	return fw, nil
}

func (fw *FileWriter) Close() error {

	if fw.file != nil {
		fw.writeDataLength()

		fw.file.Close()
		fw.file = nil
		fw.dataLength = 0
	}

	return nil
}

func (fw *FileWriter) Write(data []byte) (n int, err error) {
	n, err = fw.file.Write(data)
	fw.dataLength += uint32(n)
	return n, err
}

func (fw *FileWriter) writeConfig() error {

	_, err := fw.file.Seek(0, os.SEEK_SET)
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

		if _, err = encodeAndWrite(&ch, fw.file); err != nil {
			return err
		}

		// WAVE
		{
			var riffFormat = token_WAVE
			if _, err = fw.file.Write(riffFormat[:]); err != nil {
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

		if _, err = encodeAndWrite(&ch, fw.file); err != nil {
			return err
		}

		var c_data fmtData
		if err = c_data.setConfig(fw.config); err != nil {
			return err
		}

		if _, err = encodeAndWrite(&c_data, fw.file); err != nil {
			return err
		}
	}

	// data chunk header
	{
		ch = chunkHeader{
			id:   token_data,
			size: 0,
		}

		if _, err = encodeAndWrite(&ch, fw.file); err != nil {
			return err
		}
	}

	return nil
}

func (fw *FileWriter) writeDataLength() error {

	var data = make([]byte, size_chunkSize)

	size_dataChunk := size_chunkHeader + fw.dataLength

	// RIFF chunk
	{
		size := uint32(size_Wave + size_FmtChunk + size_dataChunk)
		pos := int64(size_chunkID)

		byteOrder.PutUint32(data, size)
		fw.file.Seek(pos, os.SEEK_SET)
		if _, err := fw.file.Write(data); err != nil {
			return err
		}
	}

	// data chunk
	{
		size := fw.dataLength
		pos := int64(size_RiffHeader + size_FmtChunk + size_chunkID)

		byteOrder.PutUint32(data, size)
		fw.file.Seek(pos, os.SEEK_SET)
		if _, err := fw.file.Write(data); err != nil {
			return err
		}
	}

	return nil
}
