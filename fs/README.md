# Package 'fs'

[![Godoc Documentation](https://godoc.org/github.com/Scalingo/go-utils/fs?status.svg)](https://godoc.org/github.com/Scalingo/go-utils/fs)

## Purpose

This package contains common methods around file manipulation. As well as an
in-memory implementation of these operations to use for tests especially.

### Constructors

* `NewOsFs()`
* `NewMemFs()`

## Warning

* `Chown`
* `Link`
* `Symlink`

Are no-op for MemFs implementation
