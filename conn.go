package main

import (
	"io"
	"net"
)

type PooledConn struct {
	conn             net.Conn
	pool             chan (*PooledConn)
	reader           *Reader
	firstByteBuf     []byte
	knownValueReader *KnownValueReader
}

func (pc *PooledConn) release() {
	pc.pool <- pc
}

func (pc *PooledConn) Read(p []byte) (int, error) {
	i, e := pc.knownValueReader.Read(p)
	if e == io.EOF {
		pc.release()
	}
	return i, e
}
