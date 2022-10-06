package socksproxy

import "github.com/armon/go-socks5"
import "fmt"

func StartServer(ip string, port string) {
    conf := &socks5.Config{}
    server, err := socks5.New(conf)
    if err != nil {
      panic(err)
    }

    // Create SOCKS5 proxy on localhost port 8000
    fmt.Println("Starting SOCKS proxy on: ", ip,":",port)
    if err := server.ListenAndServe("tcp", ip+":"+port); err != nil {
      panic(err)
    }
}
