package csv

import (
	"bytes"

	"github.com/gocarina/gocsv"
)

func MarshalCsvWithBOM(records interface{}) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.Write([]byte{0xEF, 0xBB, 0xBF}) // UTF-8
	if err != nil {
		return nil, err
	}

	err = gocsv.Marshal(records, &buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
