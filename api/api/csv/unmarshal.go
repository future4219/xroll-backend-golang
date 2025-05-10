package csv

import (
	"bytes"
	"io"

	"github.com/gocarina/gocsv"
)

type bomSkipper struct {
	reader     io.Reader
	skippedBOM bool
}

func (bs *bomSkipper) Read(p []byte) (int, error) {
	if !bs.skippedBOM {
		// 頭の3byteを読み込む
		bom := make([]byte, 3)
		n, err := bs.reader.Read(bom)
		if err != nil {
			return n, err
		}
		if n == 3 && bom[0] == 0xEF && bom[1] == 0xBB && bom[2] == 0xBF {
			// 頭の3byteがBOMなら読み飛ばす
		} else {
			// BOMでないなら読み込んだ3byteをpにコピーする
			copy(p, bom[:n])
		}
		bs.skippedBOM = true
	}
	return bs.reader.Read(p)
}

func UnmarshalCsvWithBOM[T any](csvData []byte) ([]T, error) {
	bsReader := &bomSkipper{reader: bytes.NewReader(csvData)}
	var records []T

	if err := gocsv.Unmarshal(bsReader, &records); err != nil {
		return nil, err
	}

	return records, nil
}
