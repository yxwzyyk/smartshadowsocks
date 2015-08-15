package sss

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
)

type Conn struct {
	net.Conn
	*Cipher
	readBuf  []byte
	writeBuf []byte
}

func NewConn(conn net.Conn, cipher *Cipher) *Conn {
	return &Conn{
		Conn:     conn,
		Cipher:   cipher,
		readBuf:  leakyBuf.Get(),
		writeBuf: leakyBuf.Get(),
	}
}

func (conn *Conn) Close() error {
	leakyBuf.Put(conn.readBuf)
	leakyBuf.Put(conn.writeBuf)
	return conn.Conn.Close()
}

func (conn *Conn) Read(b []byte) (n int, err error) {
	if conn.dec == nil {
		iv := make([]byte, conn.info.ivLen)
		if _, err = io.ReadFull(conn.Conn, iv); err != nil {
			return
		}
		if err = conn.initDecrypt(iv); err != nil {
			return
		}
	}

	cipherData := conn.readBuf
	if len(b) > len(cipherData) {
		cipherData = make([]byte, len(b))
	} else {
		cipherData = cipherData[:len(b)]
	}

	n, err = conn.Conn.Read(cipherData)
	if n > 0 {
		conn.decrypt(b[0:n], cipherData[0:n])
	}
	return
}

func (conn *Conn) Write(b []byte) (n int, err error) {
	var iv []byte
	if conn.enc == nil {
		if iv, err = conn.initEncrypt(); err != nil {
			return
		}
	}

	cipherData := conn.writeBuf
	dataSize := len(b) + len(iv)
	if dataSize > len(cipherData) {
		cipherData = make([]byte, dataSize)
	} else {
		cipherData = cipherData[:dataSize]
	}

	if iv != nil {
		copy(cipherData, iv)
	}
	conn.encrypt(cipherData[len(iv):], b)
	n, err = conn.Conn.Write(cipherData)
	return
}

func Dial(addr, server string, cipher *Cipher) (c *Conn, err error) {
	ra, err := RawAddr(addr)
	if err != nil {
		return
	}
	return DialWithRawAddr(ra, server, cipher)
}

func RawAddr(addr string) (buf []byte, err error) {
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, fmt.Errorf("shadowsocks: address error %s %v", addr, err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("shadowsocks: invalid port %s", addr)
	}

	hostLen := len(host)
	l := 1 + 1 + hostLen + 2
	buf = make([]byte, l)
	buf[0] = 3
	buf[1] = byte(hostLen)
	copy(buf[2:], host)
	binary.BigEndian.PutUint16(buf[2+hostLen:2+hostLen+2], uint16(port))
	return
}

func DialWithRawAddr(rawaddr []byte, server string, cipher *Cipher) (c *Conn, err error) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		return
	}
	c = NewConn(conn, cipher)
	if _, err = c.Write(rawaddr); err != nil {
		c.Close()
		return nil, err
	}
	return
}
