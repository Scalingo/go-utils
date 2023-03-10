package io

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopier_Copy(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	wd = filepath.Join(wd, "tmp")
	err = os.Mkdir(wd, 0700)
	require.NoError(t, err)
	defer os.RemoveAll(wd)

	t.Run("it should behave like io.Copy", func(t *testing.T) {
		src := bytes.NewBuffer([]byte("hello"))
		dst := &bytes.Buffer{}
		copier := NewCopier()
		n, err := copier.Copy(dst, src)
		require.NoError(t, err)
		assert.Equal(t, int64(5), n)
		assert.Equal(t, "hello", dst.String())
	})

	t.Run("it should work for files to", func(t *testing.T) {
		src, err := os.CreateTemp(wd, "go-utils-io")
		require.NoError(t, err)
		defer os.Remove(src.Name())
		dst, err := os.CreateTemp(wd, "go-utils-io")
		require.NoError(t, err)
		defer os.Remove(dst.Name())

		_, err = src.WriteString("hello")
		require.NoError(t, err)
		// Get back the reader at the beginning of the file
		_, err = src.Seek(0, io.SeekStart)
		require.NoError(t, err)

		copier := NewCopier()
		n, err := copier.Copy(dst, src)
		require.NoError(t, err)
		assert.Equal(t, int64(5), n)

		require.NoError(t, dst.Close())
		require.NoError(t, src.Close())

		body, err := os.ReadFile(dst.Name())
		require.NoError(t, err)
		assert.Equal(t, "hello", string(body))
	})

	t.Run("it should not impact cache if option is set", func(t *testing.T) {
		src, err := os.CreateTemp(wd, "go-utils-io")
		require.NoError(t, err)
		defer os.Remove(src.Name())
		dst, err := os.CreateTemp(wd, "go-utils-io")
		require.NoError(t, err)
		defer os.Remove(dst.Name())

		// Exactly 1024 bytes in lorem.txt
		fixture, err := os.ReadFile("testdata/lorem.txt")
		require.NoError(t, err)

		// Write down 100MB
		for i := 0; i < 1024*100; i++ {
			_, err = src.Write(fixture)
			require.NoError(t, err)
		}

		// Get back the reader at the beginning of the file
		_, err = src.Seek(0, io.SeekStart)
		require.NoError(t, err)

		copier := NewCopier(WithNoDiskCache)
		_, err = copier.Copy(dst, src)
		require.NoError(t, err)
		require.NoError(t, dst.Close())
		require.NoError(t, src.Close())
	})
}
