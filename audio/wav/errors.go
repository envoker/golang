package wav

import (
	"errors"
	"fmt"
)

func newError(message string) error {
	return errors.New(fmt.Sprint("wav.", message))
}

var (
	ErrorItIsNotRiffFile = newError("ItIsNotRiffFile")
	ErrorItIsNotWave     = newError("ItIsNotWave")
	ErrorWrongDataLen    = newError("WrongDataLen")
	ErrorChunkStructure  = newError("ChunkStructure")
	ErrorEOF             = newError("EOF")
	ErrorBytesPerSample  = newError("Wrong BytesPerSample")
)
