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

package data

import "io"

type reader struct {
	_   [0]func()
	r   io.Reader
	buf [8]byte
}

func (r *reader) Close() error {
	if c, ok := r.r.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

// NewReader creates a simple Reader struct from the base io.Reader provided.
func NewReader(r io.Reader) Reader {
	return &reader{r: r}
}
func (r *reader) Int() (int, error) {
	v, err := r.Uint64()
	if err != nil {
		return 0, err
	}
	return int(v), nil
}
func (r *reader) Uint() (uint, error) {
	v, err := r.Uint64()
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}
func (r *reader) Bool() (bool, error) {
	v, err := r.Uint8()
	if err != nil {
		return false, err
	}
	return v == 1, nil
}
func (r *reader) Int8() (int8, error) {
	v, err := r.Uint8()
	if err != nil {
		return 0, err
	}
	return int8(v), nil
}
func (r *reader) ReadInt(p *int) error {
	v, err := r.Int()
	if err != nil {
		return err
	}
	*p = v
	return nil
}
func (r *reader) Int16() (int16, error) {
	v, err := r.Uint16()
	if err != nil {
		return 0, err
	}
	return int16(v), nil
}
func (r *reader) Int32() (int32, error) {
	v, err := r.Uint32()
	if err != nil {
		return 0, err
	}
	return int32(v), nil
}
func (r *reader) Int64() (int64, error) {
	v, err := r.Uint64()
	if err != nil {
		return 0, err
	}
	return int64(v), nil
}
func (r *reader) Uint8() (uint8, error) {
	n, err := r.r.Read(r.buf[0:1])
	if err != nil {
		return 0, err
	}
	if n < 1 {
		return 0, io.EOF
	}
	return r.buf[0], nil
}
func (r *reader) Bytes() ([]byte, error) {
	t, err := r.Uint8()
	if err != nil {
		return nil, err
	}
	var l uint64
	switch t {
	case 0:
		return nil, nil
	case 1, 2:
		n, err2 := r.Uint8()
		if err2 != nil {
			return nil, err2
		}
		l = uint64(n)
	case 3, 4:
		n, err2 := r.Uint16()
		if err2 != nil {
			return nil, err2
		}
		l = uint64(n)
	case 5, 6:
		n, err2 := r.Uint32()
		if err2 != nil {
			return nil, err2
		}
		l = uint64(n)
	case 7, 8:
		n, err2 := r.Uint64()
		if err2 != nil {
			return nil, err2
		}
		l = n
	default:
		return nil, ErrInvalidType
	}
	if l == 0 {
		// NOTE(dij): Technically we should return (nil, nil)
		//            But! Our spec states that 0 size should be ID:0
		//            NOT ID:0,SIZE:0. So something made a fucky wucky here.
		return nil, io.ErrUnexpectedEOF
	}
	if l > MaxSlice {
		return nil, ErrTooLarge
	}
	var (
		b = make([]byte, l)
		n int
	)
	if n, err = io.ReadFull(r.r, b); err != nil {
		switch {
		case err == io.EOF:
		case err == ErrLimit:
		default:
			return nil, err
		}
	}
	if uint64(n) != l {
		return b[:n], io.EOF
	}
	return b, nil
}
func (r *reader) ReadUint(p *uint) error {
	v, err := r.Uint()
	if err != nil {
		return err
	}
	*p = v
	return nil
}
func (r *reader) ReadInt8(p *int8) error {
	v, err := r.Int8()
	if err != nil {
		return err
	}
	*p = v
	return nil
}
func (r *reader) ReadBool(p *bool) error {
	v, err := r.Bool()
	if err != nil {
		return err
	}
	*p = v
	return nil
}
func (r *reader) Uint16() (uint16, error) {
	_ = r.buf[1]
	n, err := io.ReadFull(r.r, r.buf[0:2])
	if err != nil {
		return 0, err
	}
	if n < 2 {
		return 0, io.EOF
	}
	return uint16(r.buf[1]) | uint16(r.buf[0])<<8, nil
}
func (r *reader) Uint32() (uint32, error) {
	_ = r.buf[3]
	n, err := io.ReadFull(r.r, r.buf[0:4])
	if err != nil {
		return 0, err
	}
	if n < 4 {
		return 0, io.EOF
	}
	return uint32(r.buf[3]) | uint32(r.buf[2])<<8 | uint32(r.buf[1])<<16 | uint32(r.buf[0])<<24, nil
}
func (r *reader) Uint64() (uint64, error) {
	_ = r.buf[7]
	n, err := io.ReadFull(r.r, r.buf[:])
	if err != nil {
		return 0, err
	}
	if n < 8 {
		return 0, io.EOF
	}
	return uint64(r.buf[7]) | uint64(r.buf[6])<<8 | uint64(r.buf[5])<<16 | uint64(r.buf[4])<<24 |
		uint64(r.buf[3])<<32 | uint64(r.buf[2])<<40 | uint64(r.buf[1])<<48 | uint64(r.buf[0])<<56, nil
}
func (r *reader) ReadInt16(p *int16) error {
	v, err := r.Int16()
	if err != nil {
		return err
	}
	*p = v
	return nil
}
func (r *reader) ReadInt32(p *int32) error {
	v, err := r.Int32()
	if err != nil {
		return err
	}
	*p = v
	return nil
}
func (r *reader) ReadInt64(p *int64) error {
	v, err := r.Int64()
	if err != nil {
		return err
	}
	*p = v
	return nil
}
func (r *reader) ReadUint8(p *uint8) error {
	v, err := r.Uint8()
	if err != nil {
		return err
	}
	*p = v
	return nil
}
func (r *reader) Float32() (float32, error) {
	v, err := r.Uint32()
	if err != nil {
		return 0, nil
	}
	return float32FromInt(v), nil
}
func (r *reader) Float64() (float64, error) {
	v, err := r.Uint64()
	if err != nil {
		return 0, nil
	}
	return float64FromInt(v), nil
}
func (r *reader) ReadBytes(p *[]byte) error {
	v, err := r.Bytes()
	if err != nil {
		return err
	}
	*p = v
	return nil
}
func (r *reader) Read(b []byte) (int, error) {
	return r.r.Read(b)
}
func (r *reader) ReadUint16(p *uint16) error {
	v, err := r.Uint16()
	if err != nil {
		return err
	}
	*p = v
	return nil
}
func (r *reader) ReadUint32(p *uint32) error {
	v, err := r.Uint32()
	if err != nil {
		return err
	}
	*p = v
	return nil
}
func (r *reader) ReadUint64(p *uint64) error {
	v, err := r.Uint64()
	if err != nil {
		return err
	}
	*p = v
	return nil
}
func (r *reader) ReadString(p *string) error {
	v, err := r.StringVal()
	if err != nil {
		return err
	}
	*p = v
	return nil
}
func (r *reader) StringVal() (string, error) {
	b, err := r.Bytes()
	if err != nil {
		return "", err
	}
	return string(b), nil
}
func (r *reader) ReadFloat32(p *float32) error {
	v, err := r.Float32()
	if err != nil {
		return err
	}
	*p = v
	return nil
}
func (r *reader) ReadFloat64(p *float64) error {
	v, err := r.Float64()
	if err != nil {
		return err
	}
	*p = v
	return nil
}
