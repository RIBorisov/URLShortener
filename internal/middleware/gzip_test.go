package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompressWriter(t *testing.T) {
	w := httptest.NewRecorder()
	cw := newCompressWriter(w)

	data := []byte("Test data")
	n, err := cw.Write(data)
	assert.NoError(t, err)
	assert.Equal(t, len(data), n)

	assert.Equal(t, "application/x-gzip", w.Header().Get("Content-Type"))

	assert.NoError(t, cw.Close())

	var buf bytes.Buffer
	zr, err := gzip.NewReader(bytes.NewBuffer(w.Body.Bytes()))
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, zr.Close())
	}()
	const maxDecompressedSize = 1024
	lr := io.LimitReader(zr, maxDecompressedSize)
	_, err = io.Copy(&buf, lr)
	assert.NoError(t, err)
	assert.Equal(t, string(data), buf.String())
}

func TestCompressWriter_WriteHeader(t *testing.T) {
	w := httptest.NewRecorder()
	cw := newCompressWriter(w)

	cw.WriteHeader(http.StatusOK)
	assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))
}

func TestCompressReader(t *testing.T) {
	data := []byte("Hola, amigos!")
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write(data)
	assert.NoError(t, err)
	err = zw.Flush()
	assert.NoError(t, err)
	assert.NoError(t, zw.Close()) // it is necessary that the writer be closed here for a full write to the buffer
	r, err := newCompressReader(io.NopCloser(&buf))
	assert.NoError(t, err)

	var readBuf bytes.Buffer
	n, err := io.Copy(&readBuf, r)
	assert.NoError(t, err)
	assert.Equal(t, len(data), int(n)) // int64 -> int for a valid comparison
	assert.Equal(t, string(data), readBuf.String())

	err = r.Close()
	assert.NoError(t, err)
}

func TestCompressReader_Error(t *testing.T) {
	data := []byte("Hello, World!")
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	defer func() {
		assert.NoError(t, zw.Close())
	}()
	_, err := zw.Write(data)
	assert.NoError(t, err)

	assert.NoError(t, zw.Flush())

	_, err = newCompressReader(io.NopCloser(&buf))
	assert.NoError(t, err)

	buf.Truncate(len(buf.Bytes()))

	_, err = newCompressReader(io.NopCloser(&buf))
	assert.Error(t, err)
}
