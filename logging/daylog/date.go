package daylog

import (
	"github.com/envoker/golang/time/date"
)

const formatDateFileName = "2006-01-02.log"

func dateToFileName(d date.Date) (fileName string) {
	return d.Format(formatDateFileName)
}

func dateFromFileName(fileName string) (date.Date, error) {
	return date.Parse(formatDateFileName, fileName)
}
