package compressor

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
)

type gzipCompressor struct {
	outTempfile *os.File
	outWriter   io.Writer
}

func NewGzipCompressor() (*gzipCompressor, error) {
	var (
		comp gzipCompressor
		err  error
	)

	comp.outTempfile, err = ioutil.TempFile("", "compress_out")
	if err != nil {
		return nil, err
	}

	comp.outWriter = gzip.NewWriter(comp.outTempfile)

	return &comp, nil
}

func (c *gzipCompressor) Compress(in io.ReadCloser) (io.Reader, error) {
	_, err := io.Copy(c.outWriter, in)
	if err != nil {
		return nil, err
	}

	return c.outTempfile, nil
}

func (c *gzipCompressor) Cleanup() error {
	err := os.Remove(c.outTempfile.Name())
	if err != nil {
		return err
	}

	return nil
}
