package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrOffsetIrregularFile   = errors.New("offset irregular file")
	ErrNoLimitIrregularFile  = errors.New("no limit input irregular file")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Open input file
	fdf, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fdf.Close()

	// Get information about input file
	fileInfo, err := fdf.Stat()
	if err != nil {
		return err
	}

	// Check if input file is directory
	if fileInfo.Mode().IsDir() {
		return ErrUnsupportedFile
	}

	// Check if irregular input file and no limit
	if !fileInfo.Mode().IsRegular() && limit == 0 {
		return ErrNoLimitIrregularFile
	}

	// Check offset
	if fileInfo.Mode().IsRegular() && offset > fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}
	if offset > 0 {
		ret, err := fdf.Seek(offset, io.SeekStart)
		if err != nil {
			return err
		}
		if ret != offset {
			return ErrOffsetIrregularFile
		}
	}

	// Open output file
	fdt, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer fdt.Close()

	// Size
	size := fileInfo.Size() - offset
	if limit != 0 && (size == 0 || limit < size) {
		size = limit
	}

	// Copying
	stopCopy := false
	var copiedBytes int64
	var copyBlockSize int64 = 4096
	var copiedBlockSize int64
	for !stopCopy {
		// For not regular input files
		if copiedBytes+copyBlockSize > size {
			copyBlockSize = size - copiedBytes
			stopCopy = true
		}

		copiedBlockSize, err = io.CopyN(fdt, fdf, copyBlockSize)
		if err != nil {
			stopCopy = true
		}
		copiedBytes += copiedBlockSize

		fmt.Print(progressBar(copiedBytes, size))
	}
	fmt.Println()

	return nil
}

func progressBar(copiedBytes, size int64) string {
	var percentage int64 = 100
	if size != 0 {
		percentage = copiedBytes * 100 / size
	}
	return fmt.Sprintf("\r%d bytes (%d%%)", copiedBytes, percentage)
}
