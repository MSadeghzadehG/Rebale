package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

// in the name of God.

type Rebale interface {
	Connect(address string, poolSize int) error
	Ping() error
	Get(key string) (io.Reader, error)
	Set(key string, value io.Reader, length int) error
	Close() error
}

var (
	clrf        = []byte("\r\n")
	getTemplate = "*2\r\n$3\r\nget\r\n$%v\r\n%v\r\n"
	setTemplate = "*3\r\n$3\r\nset\r\n$%v\r\n%v\r\n$%v\r\n"
)

type MyRebaleImpl struct {
	connPool chan *PooledConn
}

func (c *MyRebaleImpl) Connect(address string, poolSize int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c.connPool = make(chan *PooledConn, poolSize)

	for i := 0; i < poolSize; i++ {
		dialer := net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 10 * time.Second,
		}
		conn, err := dialer.DialContext(ctx, "tcp4", address)
		if err != nil {
			return err
		}
		r := NewReader(conn)
		c.connPool <- &PooledConn{pool: c.connPool, conn: conn, reader: r, firstByteBuf: make([]byte, 1), knownValueReader: NewKnownValueReader(r)}
	}
	return nil
}

func (c *MyRebaleImpl) Ping() error {
	command := bytes.NewReader([]byte("*1\r\n$4\r\nping\r\n"))
	pconn := <-c.connPool
	r, err := c.send(pconn, command)
	if err != nil {
		return err
	}
	s, err := r.ReadResponseMeta()
	if err != nil {
		return err
	}
	if s != "PONG" {
		return errors.New("")
	}
	return nil
}

func (c *MyRebaleImpl) Set(key string, value io.Reader, length int) error {
	pconn := <-c.connPool
	cmd := fmt.Sprintf(setTemplate, len(key), key, length)
	metaReader := bytes.NewReader([]byte(cmd))
	r, err := c.send(pconn, io.MultiReader(metaReader, io.LimitReader(value, int64(length)), bytes.NewReader(clrf)))
	if err != nil {
		return err
	}
	_, err = r.ReadResponseMeta()
	pconn.release()
	if err != nil {
		return err
	}
	return nil
}

func (c *MyRebaleImpl) Get(key string) (io.Reader, error) {
	pconn := <-c.connPool
	r, err := c.send(pconn, bytes.NewReader([]byte(fmt.Sprintf(getTemplate, len(key), key))))
	if err != nil {
		return nil, err
	}
	length, err := r.ReadInt()
	if err != nil {
		return nil, err
	}
	pconn.knownValueReader.SetValueLen(int64(length))
	return pconn, nil
}

func (c *MyRebaleImpl) Close() error {
	return nil
}

func (c *MyRebaleImpl) send(pconn *PooledConn, command io.Reader) (*Reader, error) {
	pconn.conn.SetDeadline(time.Now().Add(5 * time.Second))
	_, err := io.Copy(pconn.conn, command)
	if err != nil {
		return nil, err
	}
	return pconn.reader, nil
}
