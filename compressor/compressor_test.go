package compressor

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompress(t *testing.T) {
	type test struct {
		setup func() (io.ReadCloser, int64, error)
		err   error
	}

	tests := map[string]test{
		"succeeds from file": test{
			setup: func() (io.ReadCloser, int64, error) {
				filename := "../testdata/TestProcessEligibleChannel2_0_TestProcessEligibleChannel2_0_CODES.csv"
				instat, err := os.Stat(filename)
				if err != nil {
					return nil, 0, err
				}

				file, err := os.Open(filename)
				return file, instat.Size(), err
			},
			err: nil,
		},
		"succeeds from buffer": test{
			setup: func() (io.ReadCloser, int64, error) {
				input := strings.NewReader("This is a buffer!")
				return ioutil.NopCloser(input), input.Size(), nil
			},
			err: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			innerAssert := require.New(t)

			// Always use a new compressor
			compressor, err := NewGzipCompressor()
			innerAssert.Nil(err)

			input, inSize, err := tt.setup()
			innerAssert.Nil(err)

			// assert that the compressor did in fact write a file with no errors
			out, err := compressor.Compress(input)
			innerAssert.NotEmpty(out)
			innerAssert.Nil(err)

			// check file sizes (did we actually compress?)
			outfile := compressor.tempfile.Name()
			outstat, err := os.Stat(outfile)
			innerAssert.Nil(err)

			innerAssert.LessOrEqual(outstat.Size(), inSize)

			// make sure Cleanup works
			err = compressor.Cleanup()
			innerAssert.Nil(err)

			_, err = os.Stat(outfile)
			innerAssert.True(os.IsNotExist(err))
		})
	}
}
