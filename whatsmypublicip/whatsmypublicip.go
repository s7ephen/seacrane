package whatsmypublicip
import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
)

type data struct {
      Ip_addr        string     `json:'ip_addr'`
      Remote_host    string     `json:'remote_host'`
      User_agent     string     `json:'user_agent'`
      Port           int32      `json:'port'`
      Method         string     `json:'method'`
      Via            string     `json:'via'`
      Forwarded      string     `json:'forwarded'`
}

func Getip(showall []string ){
    ifconfigmeurl := "http://ifconfig.me/all.json"
    fmt.Println("Reaching out to :", ifconfigmeurl)
    resp, getErr := http.Get(ifconfigmeurl)
    if getErr != nil {
        log.Fatal(getErr)
    }
        body, readErr := ioutil.ReadAll(resp.Body)
    if readErr != nil {
        log.Fatal(readErr)
    }
//    fmt.Println(string(body))
    data_obj := data{}
    jsonErr := json.Unmarshal(body, &data_obj)
    if jsonErr != nil {
        log.Fatal(jsonErr)
    }
    fmt.Println("\t[+] Your public facing IP address is: ", data_obj.Ip_addr)
    if len(showall)>=1{
        if showall[0] == "all" {
            fmt.Println(string(body))
        }
    }
 }
