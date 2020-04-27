# Package 'tarball'

[![Godoc Documentation](https://godoc.org/github.com/Scalingo/go-utils/tarball?status.svg)](https://godoc.org/github.com/Scalingo/go-utils/tarball)

## Purpose

This package contains helper to manipulate tarball, compressed or not.

```
tarball.Create(context.Context, string, io.Writer, CreateOpts) error
tarball.Extract(context.Context, string, io.Reader, ExtractOpts)
tarball.Tar(io.Writer, map[string]TarFileReader) error
```
