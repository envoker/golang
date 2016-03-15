package logr

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

func gzipCompressFile(fileName string) error {

	src, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(fmt.Sprintf("%s.gz", fileName))
	if err != nil {
		return err
	}
	defer dst.Close()

	archive := gzip.NewWriter(dst)
	defer archive.Close()

	var data [4096]byte

	for {
		n, err := src.Read(data[:])
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		_, err = archive.Write(data[:n])
		if err != nil {
			return err
		}
	}

	return nil
}
