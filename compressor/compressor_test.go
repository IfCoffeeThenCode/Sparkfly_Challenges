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
	// test holds the input and expected output of a test condition
	type test struct {
		// setup does any initialization for inputs to a test, like opening files
		setup func() (input io.ReadCloser, size int64, err error)

		// cleanup does any teardown that might be required after setup()
		cleanup func(in io.ReadCloser) error

		// Output goes here. In this case, since I'm not testing that Go's gzip
		// package works (thanks, Google!) I'm just making sure that either no
		// errors were returned or that they match what I expect
		err error
	}

	// This pattern is shamelessly stolen from https://github.com/golang/go/wiki/TableDrivenTests,
	// but it served us well at Pursuant Health
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
			cleanup: func(in io.ReadCloser) error {
				return in.Close()
			},
			err: nil,
		},
		"succeeds from buffer": test{
			setup: func() (io.ReadCloser, int64, error) {
				input := strings.NewReader("This is a buffer!")
				return ioutil.NopCloser(input), input.Size(), nil
			},
			cleanup: func(in io.ReadCloser) error {
				return in.Close()
			},
			err: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assertions := require.New(t)

			// Always use a new compressor
			compressor, err := NewGzipCompressor()
			assertions.Nil(err)

			input, inSize, err := tt.setup()
			assertions.Nil(err)

			// assert that the compressor did in fact write a file with no errors
			out, err := compressor.Compress(input)
			assertions.NotEmpty(out)
			assertions.Nil(err)

			// check file sizes (did we actually compress?)
			outfile := compressor.tempfile.Name()
			outstat, err := os.Stat(outfile)
			assertions.Nil(err)

			assertions.LessOrEqual(outstat.Size(), inSize)

			// make sure Cleanup works
			err = compressor.Cleanup()
			assertions.Nil(err)

			_, err = os.Stat(outfile)
			assertions.True(os.IsNotExist(err))

			err = tt.cleanup(input)
			assertions.Nil(err)
		})
	}
}
