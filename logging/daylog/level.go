package daylog

type Level int

const (
	_ Level = iota

	LEVEL_ERROR   // logs just Errors
	LEVEL_WARNING // logs Warning and Error
	LEVEL_INFO    // logs Info, Warning and Error
	LEVEL_DEBUG   // logs everything
)

var (
	key_Level = map[Level]string{
		LEVEL_ERROR:   "Error",
		LEVEL_WARNING: "Warning",
		LEVEL_INFO:    "Info",
		LEVEL_DEBUG:   "Debug",
	}

	val_Level = map[string]Level{
		"Error":   LEVEL_ERROR,
		"Warning": LEVEL_WARNING,
		"Info":    LEVEL_INFO,
		"Debug":   LEVEL_DEBUG,
	}
)

func (l Level) String() string {
	s, _ := key_Level[l]
	return s
}

func (l Level) Valid() bool {
	_, ok := key_Level[l]
	return ok
}
