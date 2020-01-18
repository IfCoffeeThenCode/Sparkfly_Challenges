package compressor

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompress(t *testing.T) {
	assert := assert.New(t)

	compressor, err := NewGzipCompressor()
	assert.Nil(err)

	type test struct {
		in  string
		err error
	}

	tests := map[string]test{
		"good": test{
			in:  "This is a test file that will be deleted",
			err: nil,
		},
	}

	for _, tt := range tests {
		inTempfile, err := ioutil.TempFile("", "compress_in")
		inTempfile.Write([]byte(tt.in))

		// assert that the compressor did in fact write a file with no errors
		out, err := compressor.Compress(ioutil.NopCloser(strings.NewReader(tt.in)))
		assert.NotEmpty(out)
		assert.Nil(err)

		infile := inTempfile.Name()
		outfile := compressor.outTempfile.Name()

		// check file sizes (did we actually compress?)
		instat, err := os.Stat(infile)
		assert.Nil(err)

		outstat, err := os.Stat(outfile)
		assert.Nil(err)

		assert.LessOrEqual(outstat.Size(), instat.Size())

		err = compressor.Cleanup()
		assert.Nil(err)

		err = os.Remove(inTempfile.Name())
		assert.Nil(err)

		// make sure cleanup actually did cleanup
		_, err = os.Stat(outfile)
		assert.True(os.IsNotExist(err))

		_, err = os.Stat(infile)
		assert.True(os.IsNotExist(err))
	}
}
