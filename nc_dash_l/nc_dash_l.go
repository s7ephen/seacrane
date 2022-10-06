package nc_dash_l

import (
	"log"
	"net"
    "fmt"
	"strconv"
    "encoding/hex"
)

func main(){
    Nc_dash_l("",8000)
}

func Nc_dash_l(ip_to_bind string, port int) {
    fmt.Println("Dollar-store netcat....")
	l, err := net.Listen("tcp", ip_to_bind+":"+strconv.Itoa(port))
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Listening on '"+ip_to_bind+"' port", strconv.Itoa(port))
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Panicln(err)
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	log.Println("Accepted new connection.")
	defer conn.Close()
	defer log.Println("Closed connection.")

	for {
		buf := make([]byte, 1024)
		size, err := conn.Read(buf)
		if err != nil {
			return
		}
		data := buf[:size]
		log.Println("Read new data from ", conn.RemoteAddr().String() )
        fmt.Println(hex.Dump(data))
        //send the data back to the client "echoD" style. 
//		conn.Write(data)
	}
}
