package main

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

func invalidRespErr(err error) error { return fmt.Errorf("reader err: Invalid response. err: %w", err) }
func redisErr(err string) error      { return fmt.Errorf("redis err: err: %s", err) }

type Reader struct {
	reader     io.Reader
	oneByteBuf []byte
	buf        *bytes.Buffer
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		reader:     r,
		oneByteBuf: make([]byte, 1),
		buf:        new(bytes.Buffer),
	}
}

func (r *Reader) ReadResponseMeta() (s string, err error) {
	_, err = r.Read(r.oneByteBuf)
	if err != nil {
		return
	}
	b := r.oneByteBuf[0]
	if err = r.ReadUnknownValueToBuf(); err != nil {
		return
	}
	s = string(r.buf.Bytes())
	if b == byte('-') {
		return "", redisErr(s)
	}
	return s, nil
}

func (r *Reader) ReadInt() (i int, err error) {
	iStr, err := r.ReadResponseMeta()
	if err != nil {
		return -1, err
	}
	i, err = strconv.Atoi(iStr)
	if err != nil {
		return -1, invalidRespErr(err)
	}
	return
}

func (r *Reader) ReadUnknownValueToBuf() (err error) {
	r.buf.Reset()
	isPrevieusBackSlashR := false
	isValidResponse := false
	for {
		_, err = r.Read(r.oneByteBuf)
		if err != nil {
			return invalidRespErr(err)
		}
		if r.oneByteBuf[0] == byte('\n') && isPrevieusBackSlashR {
			isValidResponse = true
			break
		}
		if r.oneByteBuf[0] == byte('\r') {
			isPrevieusBackSlashR = true
		} else {
			r.buf.Write(r.oneByteBuf)
			if isPrevieusBackSlashR {
				isPrevieusBackSlashR = false
			}
		}
	}
	if !isValidResponse {
		return invalidRespErr(fmt.Errorf("response is incomplete."))
	}
	return nil
}

func (r *Reader) Read(p []byte) (int, error) {
	return r.reader.Read(p)
}

type KnownValueReader struct {
	r       io.Reader
	n       int64
	clrfLen int64
}

func NewKnownValueReader(r io.Reader) *KnownValueReader {
	return &KnownValueReader{
		r:       r,
		n:       0,
		clrfLen: 2,
	}
}

func (vr *KnownValueReader) SetValueLen(n int64) {
	vr.n = n
}

func (vr *KnownValueReader) Read(p []byte) (n int, err error) {
	if vr.n <= 0 {
		io.CopyN(io.Discard, vr.r, vr.clrfLen)
		return 0, io.EOF
	}
	if int64(len(p)) > vr.n {
		p = p[0:vr.n]
	}
	n, err = vr.r.Read(p)
	vr.n -= int64(n)
	return
}
