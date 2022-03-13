package main

import (
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	//nolint:godot
	// $ stat --format "%s" testdata/input.txt
	// 6617
	inputFileSize     int64 = 6617
	comparedBlockSize int64 = 16
)

//nolint:unparam
func isSameFiles(filename1, filename2 string) (bool, error) {
	f1, err := os.Open(filename1)
	if err != nil {
		return false, err
	}
	defer f1.Close()

	f2, err := os.Open(filename2)
	if err != nil {
		return false, err
	}
	defer f2.Close()

	buf1 := make([]byte, comparedBlockSize)
	buf2 := make([]byte, comparedBlockSize)
	for {
		_, err1 := f1.Read(buf1)
		_, err2 := f2.Read(buf2)

		if !reflect.DeepEqual(buf1, buf2) {
			return false, nil
		}

		if err1 == io.EOF && err2 == io.EOF {
			break
		}

		if err1 == io.EOF || err2 == io.EOF {
			return false, nil
		}
	}
	return true, nil
}

func TestPositiveCopy(t *testing.T) {
	t.Run("offset is equal to file size", func(t *testing.T) {
		defer os.Remove("out.txt")
		err := Copy("testdata/input.txt", "out.txt", inputFileSize, int64(0))
		require.Nil(t, err)
		require.FileExists(t, "out.txt")
		fileInfo, err := os.Stat("out.txt")
		require.NoError(t, err)
		require.Equal(t, fileInfo.Size(), int64(0))
		same, err := isSameFiles("testdata/out_offset6617_limit0.txt", "out.txt")
		require.NoError(t, err)
		require.True(t, same)
	})

	t.Run("limit is equal to file size", func(t *testing.T) {
		defer os.Remove("out.txt")
		err := Copy("testdata/input.txt", "out.txt", int64(0), inputFileSize)
		require.Nil(t, err)
		require.FileExists(t, "out.txt")
		fileInfo, err := os.Stat("out.txt")
		require.NoError(t, err)
		require.Equal(t, inputFileSize, fileInfo.Size())
		same, err := isSameFiles("testdata/out_offset0_limit6617.txt", "out.txt")
		require.NoError(t, err)
		require.True(t, same)
	})

	t.Run("limit greater than file size", func(t *testing.T) {
		defer os.Remove("out.txt")
		err := Copy("testdata/input.txt", "out.txt", int64(0), inputFileSize+1)
		require.Nil(t, err)
		require.FileExists(t, "out.txt")
		fileInfo, err := os.Stat("out.txt")
		require.NoError(t, err)
		require.Equal(t, inputFileSize, fileInfo.Size())
		same, err := isSameFiles("testdata/out_offset0_limit6618.txt", "out.txt")
		require.NoError(t, err)
		require.True(t, same)
	})

	t.Run("limit less then file size", func(t *testing.T) {
		defer os.Remove("out.txt")
		err := Copy("testdata/input.txt", "out.txt", int64(0), inputFileSize-1)
		require.Nil(t, err)
		require.FileExists(t, "out.txt")
		fileInfo, err := os.Stat("out.txt")
		require.NoError(t, err)
		require.Equal(t, inputFileSize-1, fileInfo.Size())
		same, err := isSameFiles("testdata/out_offset0_limit6616.txt", "out.txt")
		require.NoError(t, err)
		require.True(t, same)
	})

	t.Run("link file", func(t *testing.T) {
		defer os.Remove("out.txt")
		err := Copy("testdata/input.txt.lnk", "out.txt", int64(0), int64(0))
		require.Nil(t, err)
		require.FileExists(t, "out.txt")
		fileInfo, err := os.Stat("out.txt")
		require.NoError(t, err)
		require.Equal(t, inputFileSize, fileInfo.Size())
		same, err := isSameFiles("testdata/out_offset0_limit0.txt", "out.txt")
		require.NoError(t, err)
		require.True(t, same)
	})

	t.Run("link file with limit", func(t *testing.T) {
		defer os.Remove("out.txt")
		err := Copy("testdata/input.txt.lnk", "out.txt", int64(0), inputFileSize-1)
		require.Nil(t, err)
		require.FileExists(t, "out.txt")
		fileInfo, err := os.Stat("out.txt")
		require.NoError(t, err)
		require.Equal(t, inputFileSize-1, fileInfo.Size())
		same, err := isSameFiles("testdata/out_offset0_limit6616.txt", "out.txt")
		require.NoError(t, err)
		require.True(t, same)
	})

	t.Run("link file with offset and limit", func(t *testing.T) {
		defer os.Remove("out.txt")
		err := Copy("testdata/input.txt.lnk", "out.txt", int64(100), int64(1000))
		require.Nil(t, err)
		require.FileExists(t, "out.txt")
		fileInfo, err := os.Stat("out.txt")
		require.NoError(t, err)
		require.Equal(t, int64(1000), fileInfo.Size())
		same, err := isSameFiles("testdata/out_offset100_limit1000.txt", "out.txt")
		require.NoError(t, err)
		require.True(t, same)
	})

	t.Run("not regular file", func(t *testing.T) {
		defer os.Remove("out.txt")
		err := Copy("/dev/urandom", "out.txt", int64(0), int64(2048))
		require.Nil(t, err)
		require.FileExists(t, "out.txt")
		fileInfo, err := os.Stat("out.txt")
		require.NoError(t, err)
		require.Equal(t, int64(2048), fileInfo.Size())
	})

	t.Run("to dev null", func(t *testing.T) {
		err := Copy("testdata/input.txt", "/dev/null", int64(0), int64(0))
		require.Nil(t, err)
	})
}

func TestNegativeCopy(t *testing.T) {
	t.Run("offset is greater than file size", func(t *testing.T) {
		defer os.Remove("out.txt")
		err := Copy("testdata/input.txt", "out.txt", inputFileSize+1, int64(0))
		require.NoFileExists(t, "out.txt")
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("directory", func(t *testing.T) {
		defer os.Remove("out.txt")
		err := Copy("testdata", "out.txt", int64(0), int64(0))
		require.NoFileExists(t, "out.txt")
		require.ErrorIs(t, err, ErrUnsupportedFile)
	})

	t.Run("can't seek", func(t *testing.T) {
		defer os.Remove("out.txt")
		err := Copy("/dev/zero", "out.txt", int64(10), int64(10))
		require.NoFileExists(t, "out.txt")
		require.ErrorIs(t, err, ErrOffsetIrregularFile)
	})

	t.Run("irregular file without limit", func(t *testing.T) {
		defer os.Remove("out.txt")
		err := Copy("/dev/urandom", "out.txt", int64(0), int64(0))
		require.NoFileExists(t, "out.txt")
		require.ErrorIs(t, err, ErrNoLimitIrregularFile)
	})
}
