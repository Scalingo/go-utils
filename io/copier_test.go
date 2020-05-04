package io

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	procmeminfo "github.com/guillermo/go.procmeminfo"
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
		src, err := ioutil.TempFile(wd, "go-utils-io")
		require.NoError(t, err)
		defer os.Remove(src.Name())
		dst, err := ioutil.TempFile(wd, "go-utils-io")
		require.NoError(t, err)
		defer os.Remove(dst.Name())

		_, err = src.WriteString("hello")
		require.NoError(t, err)
		// Get back the reader at the beginning of the file
		src.Seek(0, os.SEEK_SET)

		copier := NewCopier()
		n, err := copier.Copy(dst, src)
		require.NoError(t, err)
		assert.Equal(t, int64(5), n)

		require.NoError(t, dst.Close())
		require.NoError(t, src.Close())

		body, err := ioutil.ReadFile(dst.Name())
		require.NoError(t, err)
		assert.Equal(t, "hello", string(body))
	})

	t.Run("it should not impact cache if option is set", func(t *testing.T) {
		src, err := ioutil.TempFile(wd, "go-utils-io")
		require.NoError(t, err)
		defer os.Remove(src.Name())
		dst, err := ioutil.TempFile(wd, "go-utils-io")
		require.NoError(t, err)
		defer os.Remove(dst.Name())

		// Exactly 1024 bytes in lorem.txt
		fixture, err := ioutil.ReadFile("test-fixtures/lorem.txt")
		require.NoError(t, err)

		// Write down 100MB
		for i := 0; i < 1024*100; i++ {
			_, err = src.Write(fixture)
			require.NoError(t, err)
		}

		// Get back the reader at the beginning of the file
		src.Seek(0, os.SEEK_SET)

		mi := &procmeminfo.MemInfo{}
		mi.Update()
		cachedBefore := (*mi)["Cached"]

		copier := NewCopier(WithNoDiskCache)
		_, err = copier.Copy(dst, src)
		require.NoError(t, err)
		require.NoError(t, dst.Close())
		require.NoError(t, src.Close())

		mi.Update()
		cachedAfter := (*mi)["Cached"]

		// If the cache diff is negative, we're all good
		// If it's positive, let's consider it's noise from the system
		// And we want to check it's not more than the 100MB of our fixture
		if cachedAfter > cachedBefore {
			cacheDiff := cachedAfter - cachedBefore
			// Let's say we want to check the difference of cache is only of 10MB,
			// indeed it can't be insured since it's global to the system
			assert.True(t, cacheDiff < 50*1024*1024, "expected diff of cache < 50MB, was %vMB", cacheDiff/1024/1024)
		}
	})
}
