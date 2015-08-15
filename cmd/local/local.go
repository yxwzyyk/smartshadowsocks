package main

import (
	"runtime"
	"flag"
	"os"
	"fmt"
	"net"
	"smartshadowsocks/sss"
)

func main() {
	if runtime.NumCPU() < 2 {
		runtime.GOMAXPROCS(2)
	} else {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	sss.PrintVersion()
	
	var configFile string
	var printVer bool

	flag.BoolVar(&printVer, "v", false, "print version")
	flag.BoolVar((*bool)(&sss.SLog.Debug), "d", false, "print debug message")
	flag.StringVar(&configFile, "c", "config.json", "specify config file")	
	flag.Parse()
	
	if printVer {
		sss.PrintVersion()
		os.Exit(0)
	}
	
	exists, err := sss.IsFileExists(configFile)
	if (!exists || err != nil) {
		fmt.Printf("%s not found! \n", configFile)
		os.Exit(1)
	}
	
	sss.Config, err = sss.ParseConfig(configFile)
	if err != nil {
		fmt.Println("Config Parse Error")
		os.Exit(1)
	}
	
	err = sss.SetRouteBuffer(sss.Config.ListFile)
	if err != nil {
		fmt.Println("Route Error:", err)
		os.Exit(1)
	}
	
	run(sss.JoinHostPort(sss.Config.Local, sss.Config.LocalPort))
}

func run(listenAddr string) {
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sss.SLog.Printf("starting local socks5 server at %v ...\n", listenAddr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			sss.SLog.Println("Accept:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	sss.SLog.DebugPrintf("socks connect from %s\n", conn.RemoteAddr().String())
	closed := false
	defer func() {
		if !closed {
			sss.SLog.DebugPrintln("handleConnection connect close")
			conn.Close()
		}
	}()
	rawaddr, host, err := sss.Socks5(conn)
	if err != nil {
		sss.SLog.Println(err)
	}
	
	var remote net.Conn
	
	if sss.RouteMatch(host) {
		sss.SLog.Printf("shadowsocks connect: %v \n", host)
		remote, err = createServerConn(rawaddr, sss.JoinHostPort(sss.Config.Server, sss.Config.ServerPort))
		if err != nil {
			sss.SLog.Println("Failed connect to shadowsocks server")
			return
		}
	} else {
		sss.SLog.Printf("Connect: %v \n", host)
		remote, err = net.Dial("tcp", host)
		if err != nil {
			sss.SLog.Println("Failed connect to %v\n", host)
			return
		}
	}
	
	defer func() {
		if !closed {
			remote.Close()
		}
	}()
	
	go sss.PipeThenClose(conn, remote)
	sss.PipeThenClose(remote, conn)
	closed = true
	sss.SLog.Printf("Closed connection to: %s\n", host)
}

func createServerConn(rawaddr []byte, host string) (remote *sss.Conn, err error) {
	cipher, err := sss.NewCipher(sss.Config.Method, sss.Config.Password)
	if err != nil {
		sss.SLog.Println("Failed generating ciphers:", err)
	}
	remote, err = sss.DialWithRawAddr(rawaddr, host, cipher.Copy())
	if err != nil {
		sss.SLog.Println("Error connecting to shadowsocks server:", err)
	}
	return
}