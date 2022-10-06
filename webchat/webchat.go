/*

A http based webchat implemented in a single file. 

Serves as a good way to quickly chat or share copy/paste buffers between machines


TODO:

   Add encryption and decrption on the fly (will require some annoying javascript decryption, which is why
I omitted it for the time being, because I dont know javascript that well).


*/

package webchat

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
    "fmt"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var savedsocketreader []*socketReader

func socketReaderCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("socket request")
	if savedsocketreader == nil {
		savedsocketreader = make([]*socketReader, 0)
	}

	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
		r.Body.Close()

	}()
	con, _ := upgrader.Upgrade(w, r, nil)

	ptrSocketReader := &socketReader{
		con: con,
	}

	savedsocketreader = append(savedsocketreader, ptrSocketReader)

	ptrSocketReader.startThread()
}

//socketReader struct
type socketReader struct {
	con  *websocket.Conn
	mode int
	name string
}

func (i *socketReader) broadcast(str string) {
	for _, g := range savedsocketreader {

		if g == i {
            // dont send user messages back to the same user
			continue
		}

		if g.mode == 1 {
            // wait for username before sending anything
			continue
		}
		g.writeMsg(i.name, str)
	}
}

func (i *socketReader) read() {
	_, b, er := i.con.ReadMessage()
	if er != nil {
		panic(er)
	}
	log.Println(i.name + " " + string(b))
	log.Println(i.mode)

	if i.mode == 1 {
		i.name = string(b)
		i.writeMsg("Sa7Chat", "Welcome '"+i.name+"', Ok, all future input below will be sent into the chatroom.")
		i.mode = 2 // real msg mode

		return
	}

	i.broadcast(string(b))

	log.Println(i.name + " " + string(b))
}

func (i *socketReader) writeMsg(name string, str string) {
	i.con.WriteMessage(websocket.TextMessage, []byte("<b>"+name+": </b>"+str))
}

func (i *socketReader) startThread() {
	i.writeMsg("Sa7Chat", "To start, send your username in the chat input below, it will be used afterwards to identify you.")
	i.mode = 1

	go func() {
		defer func() {
			err := recover()
			if err != nil {
				log.Println(err)
			}
			log.Println("Socket reader has terminated.")
		}()

		for {
			i.read()
		}

	}()
}

func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, Indexhtml)
}

func socketjavascript(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, Mysocketjs)
}

/*
func main() {
    Startsa7chat("8080")
}
*/

func Startsa7chat(ip string, port string) {
	myhttp := http.NewServeMux()
//	fs := http.FileServer(http.Dir("./views/"))
//	myhttp.Handle("/", http.StripPrefix("", fs))
    myhttp.HandleFunc("/", homePage)
    myhttp.HandleFunc("/mysocket.js", socketjavascript)
	myhttp.HandleFunc("/socket", socketReaderCreate)

    resetGlobals(ip,port)
    prefports:=map[string]string{"80":"","8080":"","443":"","8000":""}
    value,portInPrefPorts := prefports[port] // <-- Since golang doesnt have a .contains() like python for arrays
                                            //      this is a hack to put a list of stuff in a dictionary and treat a variable
                                            //      used to fetch from that dictionary (aka "map") as a boolean 
    if portInPrefPorts {} else {
        log.Println("Note: Some browsers seem to block websockets on ports other than 80,8080,443,8000",value)
    }
	log.Println("\t[+] http://"+spawnIP+":"+spawnPort)
	http.ListenAndServe(spawnIP+":"+spawnPort, myhttp)
}

// -------- Served Files Below -------

/* 
This might be a bit weird looking so here is what's going on:

Indexhtml and Mysocketjs contain the contents of what was previously served as files from:
    ./views/mysocket.js
    ./views/index.html

They were instead served from globals, which works, but if we want to change the port of the webserver, i
we also have to update the static files to contain the port for the javascript to connect to.
so we need a helperfunction to reset the globals with the new port value, since I dont know how to 
use templates in golang.
*/

var spawnIP = ``
var spawnPort = ``
var Mysocketjs = ``
var Indexhtml = ``

func resetGlobals(ip string, port string) {
//Indentation is not used to avoid problems with the verbatim strings.
    spawnIP = ip
    spawnPort = port

    Mysocketjs = `class MySocket{
    constructor(){
        this.mysocket =  null;
        this.vMsgContainer = document.getElementById("msgcontainer");
        this.vMsgIpt = document.getElementById("ipt");
    }

    showMessage(text, myself){
        var div = document.createElement("div"); 
        div.innerHTML = text;
        var cself = (myself)? "self" : "";
        div.className="msg " + cself;
        this.vMsgContainer.appendChild(div);
    }

    send(){
        var txt = this.vMsgIpt.value; 
        this.showMessage("<b>Me</b> " + txt,true);
        this.mysocket.send(txt);
        this.vMsgIpt.value = ""
    }
 
    keypress(e){
        if (e.keyCode == 13) {
            this.send();
        }
    }

    connectSocket(){
        console.log("socket initializing...");
        var socket = new WebSocket("ws://`+spawnIP+`:`+spawnPort+`/socket"); 
        this.mysocket = socket;

        socket.onmessage = (e)=>{  
           this.showMessage(e.data,false);
 
        }
        socket.onopen =  ()=> {
           console.log("socket opend")
        };  
        socket.onclose = ()=>{
           console.log("socket close")
        }
    }
}`

    Indexhtml = `<!DOCTYPE html>
<html>
    <head>
        <title>Seacrane Webchat</title>
        <style>
            .msg{
                background-color: aqua;
                padding: 10px;
                margin-bottom: 10px;
            }
            .self{
                background-color: rgb(0, 153, 255);
            }

        </style>
    </head>
    <body>
        <div id="msgcontainer" style="height:500px;background:#eee;overflow-y: auto;"></div>
        <div style="overflow: hidden;">
            <input id="ipt" style="width: 100%; padding: 10px;" onkeypress="mysocket.keypress(event)" type="text"/>
        </div>

        <script src="mysocket.js"></script>
        <script>
            var mysocket = new MySocket()
            mysocket.connectSocket();
        </script>
    </body>
</html>`
}
