package getip

import (
    "fmt"
    "net"
)

func LocalAddresses() {
    ifaces, err := net.Interfaces()
    fmt.Println("getting local interfaces...")
    if err != nil {
        fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
        return
    }
    index:=1
    for _, i := range ifaces {
        addrs, err := i.Addrs()
        fmt.Println("\t[",index,"]","---",i, addrs)
        if err != nil {
            fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
            continue
        }
        for _, a := range addrs {
            switch v := a.(type) {
            case *net.IPAddr:
                fmt.Printf("%v : %s (%s)\n", i.Name, v, v.IP.DefaultMask())
            }

        }
        index++
    }
}

func main() {
    LocalAddresses()
}
