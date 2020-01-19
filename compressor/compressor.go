package compressor

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
)

type gzipCompressor struct {
	tempfile *os.File
	writer   io.Writer
}

// NewGzipCompressor returns creates a tempfile-backed gzip writer
func NewGzipCompressor() (*gzipCompressor, error) {
	var (
		comp gzipCompressor
		err  error
	)

	// "compress_out" could be changed if need be
	comp.tempfile, err = ioutil.TempFile("", "compress_out")
	if err != nil {
		return nil, err
	}

	comp.writer = gzip.NewWriter(comp.tempfile)

	return &comp, nil
}

// Compress takes any readable buffer or file and returns a compressed tempfile
func (c *gzipCompressor) Compress(in io.ReadCloser) (io.Reader, error) {
	// Nice that gzip.Writer implements io.Writer
	_, err := io.Copy(c.writer, in)
	if err != nil {
		return nil, err
	}

	return c.tempfile, nil
}

// Cleanup removes the underlying tempfile for a given compressor
func (c *gzipCompressor) Cleanup() error {
	err := os.Remove(c.tempfile.Name())
	if err != nil {
		return err
	}

	return nil
}
