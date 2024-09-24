package rmany

import (
	"bytes"
	"context"
	"io"
	"io/fs"

	pm "github.com/takanoriyanagitani/go-post-many-http"
)

type RawSourceManyFsDir struct {
	// The trusted path to the directry.
	TrustedDirName string

	fs.ReadDirFS

	MaxBodySize int64
}

// Copies bytes from the reader to the buffer using [io.LimitedReader].
func (r RawSourceManyFsDir) ReaderToBuffer(
	rdr io.Reader,
	buf *bytes.Buffer,
) error {
	buf.Reset()
	limited := &io.LimitedReader{
		R: rdr,
		N: r.MaxBodySize,
	}
	_, e := io.Copy(buf, limited)
	return e
}

func (r RawSourceManyFsDir) FilenameToBuffer(
	filename string,
	buf *bytes.Buffer,
) error {
	f, e := r.ReadDirFS.Open(filename)
	if nil != e {
		return e
	}
	defer f.Close()
	return r.ReaderToBuffer(f, buf)
}

func (r RawSourceManyFsDir) ToRawRequestSourceMany(
	eof error,
) pm.RawRequestSourceMany {
	dirents, e := r.ReadDirFS.ReadDir(r.TrustedDirName)
	return func(ctx context.Context) (pm.RawRequest, error) {
		var buf bytes.Buffer
		if nil != e {
			return pm.RawRequest{}, e
		}

		if nil != ctx.Err() {
			return pm.RawRequest{}, ctx.Err()
		}

		if 0 == len(dirents) {
			return pm.RawRequest{}, eof
		}

		var lix int = len(dirents) - 1
		var last fs.DirEntry = dirents[lix]
		var name string = last.Name()
		err := r.FilenameToBuffer(name, &buf)
		if nil != err {
			return pm.RawRequest{}, err
		}
		dirents = dirents[:lix]
		return buf.Bytes(), nil
	}
}
