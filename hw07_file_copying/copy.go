package main

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/cheggaaa/pb"
)

const BufferSize = 20

type copier struct {
	fromPath string
	toPath   string
	offset   int64
	limit    int64
	silence  bool
}

type progress struct {
	add   int64
	error error
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" {
		return ErrEmptySourcePath
	}
	if toPath == "" {
		return ErrEmptyDestinationPath
	}
	if offset < 0 {
		return ErrInvalidOffset
	}
	if limit < 0 {
		return ErrInvalidLimit
	}
	co := &copier{
		fromPath: fromPath,
		toPath:   toPath,
		offset:   offset,
		limit:    limit,
	}
	return co.execute()
}

func (c *copier) execute() error {
	fileInfo, err := os.Stat(c.fromPath)
	if os.IsNotExist(err) {
		return NewCopyError(c.fromPath, ErrFileNotExists)
	}
	if err != nil {
		return err
	}
	if fileInfo.IsDir() {
		return NewCopyError(c.fromPath, ErrIsDirectory)
	}
	if fileInfo.Size() == 0 {
		return NewCopyError(c.fromPath, ErrUnsupportedFile)
	}
	if c.offset > fileInfo.Size() {
		return NewCopyError(c.fromPath, ErrOffsetExceedsFileSize)
	}

	fileInfoTo, err := os.Stat(c.toPath)
	if err == nil && fileInfoTo.IsDir() {
		c.toPath = c.toPath + string(os.PathSeparator) + path.Base(c.fromPath)
	}
	if !c.silence {
		fmt.Println(c.getMessage())
	}
	return c.copyFile(func(total int64) func(<-chan progress) error {
		bar := pb.New64(total)
		bar.Start()
		return func(progressCh <-chan progress) error {
			defer bar.Finish()
			for prg := range progressCh {
				if prg.error != nil {
					return prg.error
				}
				bar.Add64(prg.add)
			}
			return nil
		}
	}(c.calculateBarTotal(fileInfo.Size())))
}

func (c *copier) copyFile(progressFunc func(progressCh <-chan progress) error) error {
	if progressFunc == nil || c.silence {
		progressFunc = func(progressCh <-chan progress) error {
			for range progressCh {
			}
			return nil
		}
	}
	rdFp, err := os.Open(c.fromPath)
	if err != nil {
		return NewCopyError(c.fromPath, err)
	}
	defer rdFp.Close()

	wrFp, err := os.Create(c.toPath)
	if err != nil {
		return NewCopyError(c.toPath, err)
	}
	defer wrFp.Close()

	if c.offset > 0 {
		_, err := rdFp.Seek(c.offset, io.SeekStart)
		if err != nil {
			return NewCopyError(c.fromPath, err)
		}
	}

	progressCh := c.copyAsync(rdFp, wrFp)
	return progressFunc(progressCh)
}

func (c *copier) copyAsync(r io.Reader, w io.Writer) <-chan progress {
	progressCh := make(chan progress)
	go func() {
		defer close(progressCh)
		buff := make([]byte, BufferSize)
		var totalWrite, read, write int64
		for {
			n, readErr := r.Read(buff)
			read = int64(n)
			if read > 0 {
				if c.limit > 0 && totalWrite+read > c.limit {
					read = c.limit - totalWrite
				}
				n, err := w.Write(buff[:read])
				if err != nil {
					progressCh <- progress{
						error: err,
					}
					return
				}
				write = int64(n)
				totalWrite += write
				progressCh <- progress{add: write, error: nil}
				if c.limit > 0 && totalWrite >= c.limit {
					break
				}
			}
			if readErr == io.EOF {
				break
			}
			if readErr != nil {
				progressCh <- progress{
					error: readErr,
				}
				return
			}
		}
	}()
	return progressCh
}

func (c *copier) getMessage() string {
	msg := fmt.Sprintf("Copying %s -> %s", c.fromPath, c.toPath)
	if c.limit > 0 {
		msg += fmt.Sprintf(", limit = %d", c.limit)
	}
	if c.offset > 0 {
		msg += fmt.Sprintf(", offset = %d", c.offset)
	}
	return msg
}

func (c *copier) calculateBarTotal(fileSize int64) int64 {
	if c.limit == 0 && c.offset == 0 {
		return fileSize
	}
	total := fileSize
	total -= c.offset
	if c.limit > 0 {
		if total > c.limit {
			total = c.limit
		}
	}
	return total
}
