package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestCopy(t *testing.T) {
	goleak.VerifyNone(t)
	cases := []struct {
		from     string
		to       string
		offset   int64
		limit    int64
		expected string
	}{
		{"input.txt", "res_offset0_limit0.txt", 0, 0, "out_offset0_limit0.txt"},
		{"input.txt", "res_offset0_limit10.txt", 0, 10, "out_offset0_limit10.txt"},
		{"input.txt", "res_offset0_limit1000.txt", 0, 1000, "out_offset0_limit1000.txt"},
		{"input.txt", "res_offset0_limit10000.txt", 0, 10000, "out_offset0_limit10000.txt"},
		{"input.txt", "res_offset100_limit1000.txt", 100, 1000, "out_offset100_limit1000.txt"},
		{"input.txt", "res_offset6000_limit1000.txt", 6000, 1000, "out_offset6000_limit1000.txt"},
	}

	getHash := func(file string) string {
		fp, err := os.Open(file)
		require.NoError(t, err)
		defer fp.Close()

		h := sha256.New()
		_, err = io.Copy(h, fp)
		require.NoError(t, err)
		return hex.EncodeToString(h.Sum(nil))
	}

	for _, tc := range cases {
		tc := tc
		title := fmt.Sprintf("%s -> %s", tc.from, tc.to)
		t.Run(title, func(t *testing.T) {
			dir := "./testdata/"
			from := dir + tc.from
			to := dir + tc.to
			err := Copy(from, to, tc.offset, tc.limit)
			require.NoError(t, err)
			defer func() {
				_ = os.Remove(to)
			}()
			require.Equal(t, getHash(dir+tc.expected), getHash(to))

			fromStat, err := os.Stat(from)
			require.NoError(t, err)

			toStat, err := os.Stat(to)
			require.NoError(t, err)

			c := &copier{fromPath: tc.from, toPath: tc.to, limit: tc.limit, offset: tc.offset, silence: true}
			require.Equal(t, toStat.Size(), c.calculateBarTotal(fromStat.Size()))
		})
	}
}

func TestInputErrors(t *testing.T) {
	cases := []struct {
		from   string
		to     string
		offset int64
		limit  int64
		err    error
	}{
		{"", "to", 0, 0, ErrEmptySourcePath},
		{"test", "", 0, 0, ErrEmptyDestinationPath},
		{"test", "test2", -1, 0, ErrInvalidOffset},
		{"test", "test2", 1, -1, ErrInvalidLimit},
		{"", "test2", 1, 1, ErrEmptySourcePath},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.err.Error(), func(t *testing.T) {
			err := Copy(tc.from, tc.to, tc.offset, tc.limit)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

func TestFileErrors(t *testing.T) {
	cases := []struct {
		from   string
		to     string
		offset int64
		limit  int64
		err    error
	}{
		{"hzhz", "to", 0, 0, ErrFileNotExists},
		{"./", "to", 0, 0, ErrIsDirectory},
		{"/dev/urandom", "./", 0, 0, ErrUnsupportedFile},
		{"./testdata/out_offset0_limit1000.txt", "./", 1001, 0, ErrOffsetExceedsFileSize},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.err.Error(), func(t *testing.T) {
			err := Copy(tc.from, tc.to, tc.offset, tc.limit)
			require.ErrorIs(t, err, tc.err)
		})
	}
}
