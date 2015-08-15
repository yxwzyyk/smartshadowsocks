package sss

import (
	"encoding/binary"
	"io"
	"strconv"
)

import (
	"errors"
	"net"
)

var (
	errVer           = errors.New("socks version not supported")
	errAuthExtraData = errors.New("socks authentication get extra data")
	errCmd           = errors.New("socks command not supported")
	errAddrType      = errors.New("socks addr type not supported")
	errReqExtraData  = errors.New("socks request get extra data")
)

const (
	socksVer = 5
)

const (
	socksCmdConnect = 1
	socksCmdBind    = 2
	socksCmdUDP     = 3
)

const (
	idVer     = 0
	idNmethod = 1
	idCmd     = 1
	idType    = 3
	idIP0     = 4
	idDmLen   = 4
	idDm0     = 5
)

const (
	typeIPv4 = 1
	typeDm   = 3
	typeIPv6 = 4
)

const (
	lenIPv4   = 3 + 1 + net.IPv4len + 2 // 3(ver+cmd+rsv) + 1addrType + ipv4 + 2port
	lenIPv6   = 3 + 1 + net.IPv6len + 2 // 3(ver+cmd+rsv) + 1addrType + ipv6 + 2port
	lenDmBase = 3 + 1 + 1 + 2           // 3(ver+cmd+rsv) + 1addrType + 1addrLen + 2port, plus addrLen
)

func Socks5(conn net.Conn) (rawaddr []byte, host string, err error) {
	if err = handShank(conn); err != nil {
		SLog.Println(err)
	}
	if err = getRequest(conn, &rawaddr, &host); err != nil {
		SLog.Println(err)
	}
	if err = reply(conn); err != nil {
		SLog.Println(err)
	}
	return
}

//Socks5握手
func handShank(conn net.Conn) error {
	buf := make([]byte, 258)
	var err error
	var n int
	if n, err = io.ReadAtLeast(conn, buf, idNmethod+1); err != nil {
		return err
	}
	if buf[idVer] != socksVer {
		return errVer
	}
	nmethod := int(buf[idNmethod])
	msgLen := nmethod + 2
	if n == msgLen {

	} else if n < msgLen {
		if _, err = io.ReadFull(conn, buf[n:msgLen]); err != nil {
			return err
		}
	} else {
		return errAuthExtraData
	}

	if _, err = conn.Write([]byte{socksVer, 0}); err != nil {
		return err
	}
	return nil
}

func getRequest(conn net.Conn, rawaddr *[]byte, host *string) (err error) {
	buf := make([]byte, 263)
	var n int
	if n, err = io.ReadAtLeast(conn, buf, idDmLen+1); err != nil {
		return
	}
	if buf[idVer] != socksVer {
		err = errVer
		return
	}

	reqLen := -1
	switch buf[idCmd] {
	case socksCmdConnect:
		SLog.DebugPrintln("socksCmdConnect")
		switch buf[idType] {
		case typeIPv4:
			reqLen = lenIPv4
			*host = net.IP(buf[idIP0 : idIP0+net.IPv4len]).String()
		case typeIPv6:
			reqLen = lenIPv6
			*host = net.IP(buf[idIP0 : idIP0+net.IPv6len]).String()
		case typeDm:
			reqLen = int(buf[idDmLen]) + lenDmBase
			*host = string(buf[idDm0 : idDm0+buf[idDmLen]])
		default:
			err = errAddrType
			return
		}
		if n == reqLen {

		} else if n < reqLen {
			if _, err = io.ReadFull(conn, buf[n:reqLen]); err != nil {
				return
			}
		} else {
			err = errAuthExtraData
			return
		}
		port := binary.BigEndian.Uint16(buf[reqLen-2 : reqLen])
		*host = net.JoinHostPort(*host, strconv.Itoa(int(port)))
		*rawaddr = buf[idType:reqLen]
	case socksCmdBind:
		SLog.DebugPrintln("socksCmdBind")
	case socksCmdUDP:
		SLog.DebugPrintln("socksCmdUDP")
	}
	return
}

func reply(conn net.Conn) error {
	_, err := conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x08, 0x33})
	return err
}
