package middleware

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type compressWriter struct {
	http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{w, gzip.NewWriter(w)}
}

// Write writes data to the gzip-compressed response writer.
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader sets the HTTP status code and headers for the gzip-compressed response writer.
func (c *compressWriter) WriteHeader(statusCode int) {
	const gzipThreshold = 300
	if statusCode < gzipThreshold || statusCode == http.StatusConflict {
		c.Header().Set("Content-Encoding", "gzip")
	}
	c.ResponseWriter.WriteHeader(statusCode)
}

// Close closes the gzip-compressed response writer.
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to create new reader: %w", err)
	}

	return &compressReader{r: r, zr: zr}, nil
}

// Read reads data from the gzip-compressed request reader.
func (c compressReader) Read(p []byte) (int, error) {
	return c.zr.Read(p)
}

// Close closes the gzip-compressed request reader.
func (c *compressReader) Close() error {
	err1 := c.r.Close()
	err2 := c.zr.Close()
	return errors.Join(err1, err2)
}
