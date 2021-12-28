package io

import (
	"io"

	"golang.org/x/sys/unix"
)

type Copier struct {
	bufferSize       int64
	noDiskCacheRead  bool
	noDiskCacheWrite bool
}

type Fder interface {
	Fd() uintptr
}

type CopierOpt func(c *Copier)

func WithBufferSize(s int64) CopierOpt {
	return func(c *Copier) {
		c.bufferSize = s
	}
}

func WithNoDiskCacheRead(c *Copier) {
	c.noDiskCacheRead = true
}

func WithNoDiskCacheWrite(c *Copier) {
	c.noDiskCacheWrite = true
}

func WithNoDiskCache(c *Copier) {
	WithNoDiskCacheWrite(c)
	WithNoDiskCacheRead(c)
}

func NewCopier(opts ...CopierOpt) Copier {
	copier := Copier{
		bufferSize: 32 * 1024,
	}
	for _, opt := range opts {
		opt(&copier)
	}
	return copier
}

// copyContent is highly inspired from io.Copy, but calls to fadvise have been
// added to prevent caching the whole content of the files during the process,
// impacting the whole OS disk cache
func (c Copier) Copy(dst io.Writer, src io.Reader) (int64, error) {
	var (
		written int64
		err     error
	)
	buf := make([]byte, c.bufferSize)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			if fdSrc, ok := src.(Fder); c.noDiskCacheRead && ok {
				// Fadvise is a system call giving instruction to the OS about how to behave
				// with the flag FADC_DONTNEED, it tell the OS to drop the disk cache
				// on a given file, on a give part of the file (initial offset + end offset)
				// http://man7.org/linux/man-pages/man2/posix_fadvise.2.html
				unix.Fadvise(int(fdSrc.Fd()), written, written+int64(nr), unix.FADV_DONTNEED)
			}

			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				if fdDst, ok := dst.(Fder); c.noDiskCacheWrite && ok {
					unix.Fadvise(int(fdDst.Fd()), written, written+int64(nw), unix.FADV_DONTNEED)
				}
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
