package main

import (
	"encoding/json"
	"flag"
	"io"
	"net/http"
//	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"

	log "github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
//	"github.com/kr/pty"
    "github.com/creack/pty"
)

type windowSize struct {
	Rows uint16 `json:"rows"`
	Cols uint16 `json:"cols"`
	X    uint16
	Y    uint16
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	l := log.WithField("remoteaddr", r.RemoteAddr)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		l.WithError(err).Error("Unable to upgrade connection")
		return
	}

	cmd := exec.Command("/usr/bin/bash", "-l")
//	cmd.Env = append(os.Environ(), "TERM=xterm")

	tty, err := pty.Start(cmd)
	if err != nil {
		l.WithError(err).Error("Unable to start pty/cmd")
		conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}
	defer func() {
		cmd.Process.Kill()
		cmd.Process.Wait()
		tty.Close()
		conn.Close()
	}()

	go func() {
		for {
			buf := make([]byte, 1024)
			read, err := tty.Read(buf)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
				l.WithError(err).Error("Unable to read from pty/cmd")
				return
			}
			conn.WriteMessage(websocket.BinaryMessage, buf[:read])
		}
	}()

	for {
		messageType, reader, err := conn.NextReader()
		if err != nil {
			l.WithError(err).Error("Unable to grab next reader")
			return
		}

		if messageType == websocket.TextMessage {
			l.Warn("Unexpected text message")
			conn.WriteMessage(websocket.TextMessage, []byte("Unexpected text message"))
			continue
		}

		dataTypeBuf := make([]byte, 1)
		read, err := reader.Read(dataTypeBuf)
		if err != nil {
			l.WithError(err).Error("Unable to read message type from reader")
			conn.WriteMessage(websocket.TextMessage, []byte("Unable to read message type from reader"))
			return
		}

		if read != 1 {
			l.WithField("bytes", read).Error("Unexpected number of bytes read")
			return
		}

		switch dataTypeBuf[0] {
		case 0:
			copied, err := io.Copy(tty, reader)
			if err != nil {
				l.WithError(err).Errorf("Error after copying %d bytes", copied)
			}
		case 1:
			decoder := json.NewDecoder(reader)
			resizeMessage := windowSize{}
			err := decoder.Decode(&resizeMessage)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Error decoding resize message: "+err.Error()))
				continue
			}
			log.WithField("resizeMessage", resizeMessage).Info("Resizing terminal")
			_, _, errno := syscall.Syscall(
				syscall.SYS_IOCTL,
				tty.Fd(),
				syscall.TIOCSWINSZ,
				uintptr(unsafe.Pointer(&resizeMessage)),
			)
			if errno != 0 {
				l.WithError(syscall.Errno(errno)).Error("Unable to resize terminal")
			}
		default:
			l.WithField("dataType", dataTypeBuf[0]).Error("Unknown data type")
		}
	}
}

func main() {
	var listen = flag.String("listen", "127.0.0.1:3000", "Host:port to listen on")
	var assetsPath = flag.String("assets", "./assets", "Path to assets")

	flag.Parse()

	r := mux.NewRouter()

	r.HandleFunc("/term", handleWebsocket)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(*assetsPath)))

	log.Info("Demo Websocket/Xterm terminal")
	log.Warn("Warning, this is a completely insecure daemon that permits anyone to connect and control your computer, please don't run this anywhere")

	if !(strings.HasPrefix(*listen, "127.0.0.1") || strings.HasPrefix(*listen, "localhost")) {
		log.Warn("Danger Will Robinson - This program has no security built in and should not be exposed beyond localhost, you've been warned")
	}

	if err := http.ListenAndServe(*listen, r); err != nil {
		log.WithError(err).Fatal("Something went wrong with the webserver")
	}
}
