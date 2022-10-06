/*
Http Serve a directory.
   Special paths:
        /hello   - print a test message
        /headers - display the Browser's headers
        /getip  - Extract the source IP address from headers
*/

package httpdir
import (
    "fmt"
    "strings"
    "net"
    "net/http"
    "path/filepath"
)

var pln = fmt.Println
var pf = fmt.Printf
var p = fmt.Printf
var fpf = fmt.Fprintf

func hello (w http.ResponseWriter, req *http.Request){
    fpf(w, "This is a test\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
    for name, headers :=  range req.Header {
        for _, h:= range headers {
            fpf(w, "%v: %v\n", name, h)
        }
    }
}

func getiphandler (w http.ResponseWriter, r *http.Request) {
    ip, err := getIP(r)
    if err != nil {
        w.WriteHeader(400)
        w.Write([]byte("No valid ip"))
    }
    w.WriteHeader(200)
    w.Write([]byte(ip))
}

func getIP(r *http.Request) (string, error) {
    //Get IP from the X-REAL-IP header
    ip := r.Header.Get("X-REAL-IP")
    netIP := net.ParseIP(ip)
    if netIP != nil {
        return ip, nil
    }

    //Get IP from X-FORWARDED-FOR header
    ips := r.Header.Get("X-FORWARDED-FOR")
    splitIps := strings.Split(ips, ",")
    for _, ip := range splitIps {
        netIP := net.ParseIP(ip)
        if netIP != nil {
            return ip, nil
        }
    }

    //Get IP from RemoteAddr
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    if err != nil {
        return "", err
    }
    netIP = net.ParseIP(ip)
    if netIP != nil {
        return ip, nil
    }
    return "", fmt.Errorf("No valid ip found")
}

func ServeDir(port string, dir string) {
    fmt.Println(port, dir)
    abspath,err := filepath.Abs(dir) //Abs() returns err and string, can not accept the single value return
    http.HandleFunc("/hello", hello)
    http.HandleFunc("/headers", headers)
    http.HandleFunc("/getip", getiphandler)
    http.Handle("/", http.FileServer(http.Dir(dir)))
    if err == nil {
        pf("[+] HTTP Server starting on port %s serving directory: `%s`\n\tAbsolute path of webroot %s\n",port, dir, abspath)
    } else {
        pln("There was an error")
    }
    http.ListenAndServe(":"+port, nil)
}