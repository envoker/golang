package wav

import (
	"encoding/binary"
	"os"
)

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

	if fw.file == nil {
		return ErrFileWriterClosed
	}

	fw.writeDataLength()

	err := fw.file.Close()
	fw.file = nil
	fw.dataLength = 0
	return err
}

func (fw *FileWriter) Write(data []byte) (n int, err error) {
	n, err = fw.file.Write(data)
	fw.dataLength += uint32(n)
	return n, err
}

func (fw *FileWriter) writeConfig() error {

	if _, err := fw.file.Seek(0, os.SEEK_SET); err != nil {
		return err
	}

	var ch chunkHeader

	// RIFF header
	{
		ch = chunkHeader{Id: token_RIFF, Size: 0}
		err := binary.Write(fw.file, binary.LittleEndian, ch)
		if err != nil {
			return err
		}

		err = binary.Write(fw.file, binary.LittleEndian, token_WAVE)
		if err != nil {
			return err
		}
	}

	// fmt chunk
	{
		ch = chunkHeader{
			Id:   token_fmt,
			Size: uint32(size_FmtData),
		}

		err := binary.Write(fw.file, binary.LittleEndian, ch)
		if err != nil {
			return err
		}

		c_data := configToFmtData(fw.config)
		err = binary.Write(fw.file, binary.LittleEndian, c_data)
		if err != nil {
			return err
		}
	}

	// data chunk header
	{
		ch = chunkHeader{
			Id:   token_data,
			Size: 0,
		}

		err := binary.Write(fw.file, binary.LittleEndian, ch)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fw *FileWriter) writeDataLength() error {

	size_dataChunk := size_chunkHeader + int(fw.dataLength)

	// RIFF chunk
	{
		pos := int64(size_chunkID)
		size := uint32(size_Format + size_FmtChunk + size_dataChunk)

		fw.file.Seek(pos, os.SEEK_SET)
		err := binary.Write(fw.file, binary.LittleEndian, size)
		if err != nil {
			return err
		}
	}

	// data chunk
	{
		size := fw.dataLength
		pos := int64(size_RiffHeader + size_FmtChunk + size_chunkID)

		fw.file.Seek(pos, os.SEEK_SET)

		err := binary.Write(fw.file, binary.LittleEndian, size)
		if err != nil {
			return err
		}
	}

	return nil
}
