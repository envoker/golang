package logr

func duplicate(data []byte) []byte {
	d := make([]byte, len(data))
	copy(d, data)
	return d
}
