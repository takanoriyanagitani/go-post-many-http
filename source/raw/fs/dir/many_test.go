package rmany_test

import (
	"testing"

	"bytes"
	"io"

	fd "github.com/takanoriyanagitani/go-post-many-http/source/raw/fs/dir"
)

func TestMany(t *testing.T) {
	t.Parallel()

	t.Run("RawSourceManyFsDir", func(t *testing.T) {
		t.Parallel()

		t.Run("ReaderToBuffer", func(t *testing.T) {
			t.Parallel()

			t.Run("empty", func(t *testing.T) {
				t.Parallel()

				rsrc := fd.RawSourceManyFsDir{MaxBodySize: 65536}
				var empty []byte
				var rdr io.Reader = bytes.NewReader(empty)
				var buf bytes.Buffer
				e := rsrc.ReaderToBuffer(rdr, &buf)
				if nil != e {
					t.Fatalf("unexpected error: %v\n", e)
				}
				var got []byte = buf.Bytes()
				if 0 != len(got) {
					t.Fatalf("unexpected length: %v\n", len(got))
				}
			})

			t.Run("helo", func(t *testing.T) {
				t.Parallel()

				rsrc := fd.RawSourceManyFsDir{MaxBodySize: 65536}
				var helo []byte = []byte("helo")
				var rdr io.Reader = bytes.NewReader(helo)
				var buf bytes.Buffer
				e := rsrc.ReaderToBuffer(rdr, &buf)
				if nil != e {
					t.Fatalf("unexpected error: %v\n", e)
				}
				var got []byte = buf.Bytes()
				if 4 != len(got) {
					t.Fatalf("unexpected length: %v\n", len(got))
				}
			})
		})
	})
}
