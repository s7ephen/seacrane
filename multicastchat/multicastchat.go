package multicastchat

import (
	"encoding/hex"
	"log"
	"net"
	"time"
    "strings"
    "fmt"
    "github.com/abiosoft/ishell"
	"github.com/dmichael/go-multicast/multicast"
)

func Multicastchat(address string) {
//	go ping(address)
    go multicast.Listen(address, msgHandler)
    fmt.Println("\t[+] Spawned multicast message listener on: ",address)
    // create new shell.
    // by default, new shell includes 'exit', 'help' and 'clear' commands.
    shell := ishell.New()
    shell.SetPrompt("seacrane>>multicastchat>> ")
    // display welcome info.
    shell.Println("Multicast Chat submenu")
    shell.AddCmd(&ishell.Cmd{
        Name: "send",
        Help: "send <message>",
        Func: func(c *ishell.Context) {
            conn, err := multicast.NewBroadcaster(address)
            if err != nil{panic(err)}
            conn.Write([]byte(strings.Join(c.Args, " ")))
        },
    })
    shell.Run()
}

func msgHandler(src *net.UDPAddr, n int, b []byte) {
	log.Println(n, "bytes read from", src)
	fmt.Println(hex.Dump(b[:n]))
	fmt.Println("-----")
}

func ping(addr string) {
	conn, err := multicast.NewBroadcaster(addr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn.Write([]byte("ping pong\n"))
		time.Sleep(1 * time.Second)
	}
}
