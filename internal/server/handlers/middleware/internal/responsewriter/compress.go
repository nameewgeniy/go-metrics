package responsewriter

import (
	"io"
	"net/http"
)

type (
	CompressResponseWriter struct {
		http.ResponseWriter
		writer io.Writer
	}
)

func NewCompressResponseWriter(w http.ResponseWriter, cwr io.Writer) *CompressResponseWriter {
	return &CompressResponseWriter{
		ResponseWriter: w,
		writer:         cwr,
	}
}

func (w CompressResponseWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}
