// Copyright (C) 2020 - 2023 iDigitalFlame
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

package wrapper

import (
	"compress/gzip"
	"compress/zlib"
	"io"
	"sync"
)

const compLevel = zlib.BestSpeed

const (
	// Zlib is the default Zlib Wrapper. This wrapper uses the default compression level. Use the 'NewZlib'
	// function to create a wrapper with a different level.
	Zlib = compress(0x0)
	// Gzip is the default Gzip Wrapper. This wrapper uses the default compression level. Use the 'NewGzip'
	// function to create a wrapper with a different level.
	Gzip = compress(0x1)
)

var (
	zlibWriterPool = sync.Pool{
		New: func() interface{} {
			w, _ := zlib.NewWriterLevel(nil, compLevel)
			return w
		},
	}
	gzipWriterPool = sync.Pool{
		New: func() interface{} {
			w, _ := gzip.NewWriterLevel(nil, compLevel)
			return w
		},
	}
	zlibReaderPool, gzipReaderPool sync.Pool
)

type compress uint8
type reader struct {
	_ [0]func()
	p *sync.Pool
	io.ReadCloser
}
type writer struct {
	_ [0]func()
	p *sync.Pool
	io.WriteCloser
}

func (r *reader) Close() error {
	if r.ReadCloser == nil {
		return nil
	}
	err := r.ReadCloser.Close()
	r.p.Put(r.ReadCloser)
	r.ReadCloser, r.p = nil, nil
	return err
}
func (w *writer) Close() error {
	if w.WriteCloser == nil {
		return nil
	}
	err := w.WriteCloser.Close()
	w.p.Put(w.WriteCloser)
	w.WriteCloser, w.p = nil, nil
	return err
}
func (c compress) Unwrap(r io.Reader) (io.Reader, error) {
	switch c {
	case Zlib:
		c := zlibReaderPool.Get()
		if c == nil {
			n, err := zlib.NewReader(r)
			if err != nil {
				return nil, err
			}
			return &reader{ReadCloser: n, p: &zlibReaderPool}, nil
		}
		var ()
		if err := c.(zlib.Resetter).Reset(r, nil); err != nil {
			zlibReaderPool.Put(c)
			return nil, err
		}
		return &reader{ReadCloser: c.(io.ReadCloser), p: &zlibReaderPool}, nil
	case Gzip:
		c := gzipReaderPool.Get()
		if c == nil {
			n, err := gzip.NewReader(r)
			if err != nil {
				return nil, err
			}
			return &reader{ReadCloser: n, p: &gzipReaderPool}, nil
		}
		var (
			n   = c.(*gzip.Reader)
			err = n.Reset(r)
		)
		if err == nil {
			return &reader{ReadCloser: n, p: &gzipReaderPool}, nil
		}
		gzipReaderPool.Put(c)
		return nil, err
	}
	return r, nil
}
func (c compress) Wrap(w io.WriteCloser) (io.WriteCloser, error) {
	switch c {
	case Zlib:
		c := zlibWriterPool.Get().(*zlib.Writer)
		c.Reset(w)
		return &writer{WriteCloser: c, p: &zlibWriterPool}, nil
	case Gzip:
		c := gzipWriterPool.Get().(*gzip.Writer)
		c.Reset(w)
		return &writer{WriteCloser: c, p: &gzipWriterPool}, nil
	}
	return w, nil
}
