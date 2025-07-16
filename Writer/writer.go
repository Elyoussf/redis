package writer

import (
	"io"
	resp "redis/RESP"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(wr io.Writer) *Writer {
	return &Writer{writer: wr}
}

func (w *Writer) Write(v resp.Value) error {
	bytes := v.Marshall()
	_, err := w.writer.Write(bytes)

	if err != nil {
		return err
	}
	return nil
}
