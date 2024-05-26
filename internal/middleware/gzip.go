package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	const gzipThreshold = 300
	if statusCode < gzipThreshold {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	// reading only header to check if the r is gzip
	//h, err := gzip.NewReader(io.LimitReader(r, 10))
	//if err != nil {
	//	return nil, fmt.Errorf("failed to get header %w", err)
	//}
	//if err = h.Close(); err != nil {
	//	return nil, fmt.Errorf("failed to close reader %w", err)
	//}
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to create new reader: %w", err)
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (int, error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
