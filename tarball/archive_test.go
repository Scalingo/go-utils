package tarball

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/Scalingo/go-utils/fs"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type failedWriter struct{}

func (w failedWriter) Write([]byte) (int, error) {
	return -1, errors.New("fail to write")
}

func testFs(t *testing.T) fs.Fs {
	srcFs := fs.NewMemFs()
	err := srcFs.Mkdir("/tmp", 0700)
	require.NoError(t, err)
	f, err := srcFs.OpenFile("/tmp/file", os.O_CREATE|os.O_WRONLY, 0600)
	require.NoError(t, err)
	fmt.Fprintf(f, "file content")
	err = f.Close()
	require.NoError(t, err)
	return srcFs
}

func TestCreate(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		expect func(*testing.T, *os.File, string)
	}{
		"it should return an error if src is empty": {
			expect: func(t *testing.T, out *os.File, tempDir string) {
				buffer := &bytes.Buffer{}
				err := Create(ctx, "", buffer, CreateOpts{})
				require.Error(t, err)
				assert.Contains(t, err.Error(), "src is empty")
			},
		},
		"it should return an error if dst is empty": {
			expect: func(t *testing.T, out *os.File, tempDir string) {
				err := Create(ctx, "/tmp/archive", nil, CreateOpts{})
				require.Error(t, err)
				assert.Contains(t, err.Error(), "dst writer is nil")
			},
		},
		"it should return an error if dst.Write returns an error": {
			expect: func(t *testing.T, out *os.File, tempDir string) {
				writer := failedWriter{}
				err := Create(ctx, "/tmp", writer, CreateOpts{})
				require.Error(t, err)
				assert.Contains(t, err.Error(), "fail to write")
			},
		},
		"it should make a valid tar.gz archive": {
			expect: func(t *testing.T, out *os.File, tempDir string) {
				srcFs := testFs(t)
				opts := CreateOpts{Fs: srcFs}
				err := Create(ctx, "/tmp", out, opts)
				require.NoError(t, err)

				// Ensuring standard tar command can read the output archive
				err = exec.Command("tar", "-C", tempDir, "-xf", out.Name()).Run()
				require.NoError(t, err)
				content, err := ioutil.ReadFile(tempDir + "/file")
				require.NoError(t, err)
				assert.Equal(t, "file content", string(content))
			},
		},
		"it should include the src directory in the archive if option is set": {
			expect: func(t *testing.T, out *os.File, tempDir string) {
				srcFs := testFs(t)
				opts := CreateOpts{Fs: srcFs, IncludeRoot: true}
				err := Create(ctx, "/tmp", out, opts)
				require.NoError(t, err)

				// Ensuring standard tar command can read the output archive
				err = exec.Command("tar", "-C", tempDir, "-xf", out.Name()).Run()
				require.NoError(t, err)
				content, err := ioutil.ReadFile(tempDir + "/tmp/file")
				require.NoError(t, err)
				assert.Equal(t, "file content", string(content))
			},
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			out, err := ioutil.TempFile("/tmp", "fs-test")
			require.NoError(t, err)
			defer os.Remove(out.Name())

			dir, err := ioutil.TempDir("/tmp", "fs-test-dir")
			require.NoError(t, err)
			defer os.RemoveAll(dir)

			c.expect(t, out, dir)
		})
	}
}
