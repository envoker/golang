package logr

import (
	"bufio"
	"os"
)

type fileWriter struct {
	config   Config
	file     *os.File
	buf      *bufio.Writer
	dataSize int64
}

func (p *fileWriter) Close() error {

	if p.buf != nil {
		p.buf.Flush()
		p.buf = nil
	}

	if p.file != nil {
		p.file.Close()
		p.file = nil
	}

	return nil
}

func (p *fileWriter) Flush() {
	if p.buf != nil {
		p.buf.Flush()
	}
}

func (p *fileWriter) openFile() error {

	fi, err := os.Stat(p.config.FileName)
	if err == nil {
		p.dataSize = fi.Size()
	} else {
		p.dataSize = 0
	}

	p.file, err = os.OpenFile(p.config.FileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	p.buf = bufio.NewWriter(p.file)

	return nil
}

func (p *fileWriter) Write(data []byte) (n int, err error) {

	if p.file != nil {
		if p.dataSize+int64(len(data)) > p.config.FileSize {
			p.Close()
			err = rotate(p.config.FileName, p.config.Count)
			if err != nil {
				return 0, err
			}
		}
	}

	if p.file == nil {
		if err = p.openFile(); err != nil {
			return 0, err
		}
	}

	n, err = p.buf.Write(data)
	p.dataSize += int64(n)

	return
}
