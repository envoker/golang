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
)

var (
	ErrAudioFormat    = errors.New("wav: invalid AudioFormat")
	ErrChannels       = errors.New("wav: invalid Channels")
	ErrSampleRate     = errors.New("wav: invalid SampleRate")
	ErrBytesPerSample = errors.New("wav: invalid BytesPerSample")

	ErrFileReaderClosed = errors.New("wav: FileReader is closed or not created")
)
