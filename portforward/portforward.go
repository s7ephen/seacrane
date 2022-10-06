package portforward

import (
	"fmt"
	"io"
	"net"
)

func StartForward(protocol string, listenport string, remoteip string, remoteport string) {
	ln, err := net.Listen(protocol, ":"+listenport)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go handleRequest(conn, protocol, remoteip, remoteport)
	}
}

func handleRequest(conn net.Conn, protocol string, remoteip string, remoteport string) {
	fmt.Println("\t[+] New Incoming connection to the Port Forward.")

	proxy, err := net.Dial(protocol, remoteip+":"+remoteport)
	if err != nil {
		panic(err)
	}

	fmt.Println("\t[+] Proxy connected...")
	go copyIO(conn, proxy) // <---- so-called "go routine" which is just an asych process to handle the connection
	go copyIO(proxy, conn) // <---- same thing for the other half of the duplex. like a select() loop cross handles in C
}

func copyIO(src, dest net.Conn) {
	defer src.Close()
	defer dest.Close()
	io.Copy(src, dest)
}

func main() {
    StartForward("tcp","8080","google.com", "80")
}
