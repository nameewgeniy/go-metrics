package responsewriter

import "net/http"

type (
	Info struct {
		status int
		size   int
	}

	LoggerResponseWriter struct {
		http.ResponseWriter
		Info *Info
	}
)

func NewLoggerResponseWriter(w http.ResponseWriter) *LoggerResponseWriter {
	return &LoggerResponseWriter{
		ResponseWriter: w,
		Info: &Info{
			status: 0,
			size:   0,
		},
	}
}

func (r *LoggerResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.Info.size += size
	return size, err
}

func (r *LoggerResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.Info.status = statusCode
}

func (i Info) Status() int {
	return i.status
}

func (i Info) Size() int {
	return i.size
}
