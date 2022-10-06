package cli

import "fmt"
import "os"
import "strings"
import "strconv"
import "github.com/abiosoft/ishell"
import "seacrane/filecrypt"
import "seacrane/hexdump"
import "seacrane/httpdir"
import "seacrane/socksproxy"
import "seacrane/httpdl"
import "seacrane/httpdrecv"
import "seacrane/portforward"
import "seacrane/findurls"
import "seacrane/whatsmypublicip"
import "seacrane/webchat"
import "seacrane/getip"
import "seacrane/base64file"
import "seacrane/zipdir"
import "seacrane/filedups"
import "seacrane/nc_dash_l"
import "seacrane/multicastchat"
import "seacrane/qrcode"

//import "reflect"

func Run_shell_test(){
    // create new shell.
    // by default, new shell includes 'exit', 'help' and 'clear' commands.
    shell := ishell.New()
    // display welcome info.
    shell.Println("Sample Interactive Shell")

    // register a function for "greet" command.
    shell.AddCmd(&ishell.Cmd{
        Name: "greet",
        Help: "greet user",
        Func: func(c *ishell.Context) {
            c.Println("Hello", strings.Join(c.Args, " "))
        },
    })

    // simulate an authentication
    shell.AddCmd(&ishell.Cmd{
        Name: "login",
        Help: "simulate a login",
        Func: func(c *ishell.Context) {
            // disable the '>>>' for cleaner same line input.
            c.ShowPrompt(false)
            defer c.ShowPrompt(true) // yes, revert after login.

            // get username
            c.Print("Username: ")
            username := c.ReadLine()

            // get password.
            c.Print("Password: ")
            password := c.ReadPassword()
            c.Println("Authentication Successful.")
            c.Print(username, password)
        },
    })

    // simulate an authentication
    shell.AddCmd(&ishell.Cmd{
        Name: "progress",
        Help: "Show Off the ishell progressbar widget",
        Func: func(c *ishell.Context) {
            c.ProgressBar().Start()
            for i := 0; i < 101; i++ {
                c.ProgressBar().Suffix(fmt.Sprint(" ", i, "%"))
                c.ProgressBar().Progress(i)
            }
            c.ProgressBar().Stop()
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "multi",
        Help: "Show off Multiple Choice",
        Func: func(c *ishell.Context) {
            choice := c.MultiChoice([]string{
                "Golangers",
                "Go programmers",
                "Gophers",
                "Goers",
            }, "What are Go programmers called ?")
            if choice == 2 {
                c.Println("You got it!")
            } else {
                c.Println("Sorry, you're wrong.")
            }
        },

    })

    shell.AddCmd(&ishell.Cmd{
        Name: "checklist",
        Help: "Show off The Checklist Widget",
        Func: func(c *ishell.Context) {
            languages := []string{"Python", "Go", "Haskell", "Rust"}
            choices := c.Checklist(languages,
                "What are your favourite programming languages ?", nil)
            c.Println("Your choices are...not really known hhahahaha",choices)
//            out := func() []string { ... } // convert index to language
//            c.Println("Your choices are", strings.Join(out(), ", "))
        },
    })

//TEMPLATE FOR NEW COMMANDS BELOW
/*
    shell.AddCmd(&ishell.Cmd{
        Name: "MultiChoice",
        Help: "Show off Multiple Choice",
        Func: 
    })
*/
    // run shell
    shell.Run()
}

func Run_shell(){
    shell := ishell.New()
    shell.SetPrompt("seacrane>> ")
    banner := (`
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⢀⣠⣴⠾⢻⣿⡟⠻⠶⢦⣤⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⣀⣤⠾⠛⠉⠀⠀⣸⠛⣷⠀⠀⠀⠀⠉⠙⠻⠶⣦⣤⣀⠀⠀⠀⠀⠀
    ⠀⠀⠐⠛⠋⠀⠀⠀⠀⠀⠀⠛⠀⠛⠂⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⠛⠒⠂⠀⠀
    ⠀⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡇⠀
    ⠀⠀⠀⢠⣤⣤⣤⠀⠀⠀⠀⢠⣤⡄⢠⣤⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⡄⠀⠀⠀
    ⠀⠀⠀⠈⠉⠉⠉⠀⠀⠀⠀⠸⠿⠇⠸⠿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⡇⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⣶⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡇⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡇⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠃⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡀⠃⡀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⠋⠈⠛⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⢸⣿⡇⢀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⣿⣿⡇⠀
    ⠀⠀⠀⠀⠀⠀⠀⢀⣴⠾⠋⢸⣿⡇⠈⠳⣦⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠈⠁⠀⠀⠈⠛⠃⠀⠀⠀⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄
████░▄▄▄░█░▄▄█░▄▄▀██░▄▄▀█░▄▄▀█░▄▄▀█░▄▄▀█░▄▄████
████▄▄▄▀▀█░▄▄█░▀▀░██░████░▀▀▄█░▀▀░█░██░█░▄▄████
████░▀▀▀░█▄▄▄█▄██▄██░▀▀▄█▄█▄▄█▄██▄█▄██▄█▄▄▄████
███████████████████████████████████████████████
██████ ...a still point on choppy seas ████████
▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀ 
`)

    shell.Println(banner)

    shell.AddCmd(&ishell.Cmd{
        Name: "hexdump",
        Help: "Hexdump the contents of a file.",
        Func: func(c *ishell.Context) {
//            fmt.Println("c.Args:", reflect.TypeOf(c.Args).Kind())
//            fmt.Println(len(c.Args))
            if len(c.Args) >= 1 {
                hexdump.Dump(c.Args[0])
            } else {
                fmt.Println("You must supply a filename.")
            }
//           fmt.Println(reflect.ValueOf(c.Args).Elem().Interface())
//           fmt.Println(reflect.ValueOf(c.Args))
//           fmt.Println(reflect.ValueOf(c.Args).Elem())
//            hexdump.Dump(strings.Join(c.Args, " "))
        },
    })
    shell.AddCmd(&ishell.Cmd{
        Name: "httpdir",
        Help: "Serve a directory as an http index.",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 2 {
                httpdir.ServeDir(c.Args[0], c.Args[1])
            } else {
                fmt.Println("httpdir <port> <directory_to_serve>")
            }
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "socksproxy",
        Help: "Starts a Socks5 proxy server. (with support for CONNECT method (aka tcp over http))",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 2 {
                socksproxy.StartServer(c.Args[0],c.Args[1])
            } else {
                fmt.Println("socksproxy <ip> <port>")
            }
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "aes_fileencrypt",
        Help: "Encrypt a file (symmetrically AES) without reading it into memory.",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 2 {
                filecrypt.StreamEncrypter(c.Args[0], c.Args[1])
            } else {
                fmt.Println("aes_fileecrypt <passphrase> <file>")
            }
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "aes_filedecrypt",
        Help: "Decrypt a file (symmetrically AES) without reading it into memory.",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 2 {
                filecrypt.StreamDecrypter(c.Args[0], c.Args[1])
            } else {
                fmt.Println("aes_fileecrypt <passphrase> <file>")
            }
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "httpdl",
        Help: "Download a http file directly to disk without storing it in memory. (like curl -O)",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 2 {
                httpdl.DoDownload(c.Args[0], c.Args[1])
            } else {
                fmt.Println("httpdl <localfilename> <url>")
            }
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "httprecv",
        Help: "Receieve a http file upload via embedded http server.",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 2 {
                httpdrecv.StartServer(c.Args[0], c.Args[1])
            } else {
                fmt.Println("httpdrecv <ip> <port-for-httpserver> (leave ip blank to bind on all addresses)")
            }
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "portforward",
        Help: "Perform a TCP or UDP port forward.\n\t\tExample:\n\t\tportforward tcp 8080 google.com 80",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 4 {
                portforward.StartForward(c.Args[0],c.Args[1] ,c.Args[2] ,c.Args[3])
            } else {
                fmt.Println("portforward <protocol> <local port> <remote_ip> <remote_port>")
            }
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "findurls",
        Help: "Find all URLs in a file printing some to screen and saving all to a file.",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 2 {
                findurls.FindURLs(c.Args[0],c.Args[1])
            } else {
                fmt.Println("findurls <file_to_search_for_urls> <file_to_save_found_urls>")
            }
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "getpublicip",
        Help: "Hit http://ifconfig.me for Public IP address info. run 'getpublicip all' for more info",
        Func: func(c *ishell.Context) {
            whatsmypublicip.Getip(c.Args)
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "webchat",
        Help: "Spawn a Web Chat server. Very useful for sharing copy/paste buffers across machines.",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 2 {
                webchat.Startsa7chat(c.Args[0], c.Args[1])
            } else {
                fmt.Println("webchat <ip to listen on> <port_to_listen_on>")
            }
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "getip",
        Help: "List all local network interfaces and their IPs (and other information)",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 0 {
                getip.LocalAddresses()
            } else {
                fmt.Println("getip (no arguments)")
            }
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "urlencode",
        Help: "URLEncode a string.",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 1 {
                base64file.URLEncode(c.Args[0])
            } else {
                fmt.Println("urlencode <string-to-encode>")
            }
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "urldecode",
        Help: "URLDecode a string.",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 1 {
                base64file.URLDecode(c.Args[0])
            } else {
                fmt.Println("urldecode <string-to-decode>")
            }
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "base64enc",
        Help: "Base64 Encode a file and write it to disk. (gzip compresses first)",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 1 {
                base64file.EncodeFile(c.Args[0])
            } else {
                fmt.Println("base64enc <file_to_base64>")
            }
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "base64dec",
        Help: "Base64 Decode a file and write it to disk. (gzip decompresses after b64 decode)",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 1 {
                base64file.DecodeFile(c.Args[0])
            } else {
                fmt.Println("base64dec <file_to_base64>")
            }
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "string_base64enc",
        Help: "Base64 Encode a string and print to screen",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 1 {
                base64file.EncodeString(c.Args[0])
            } else {
                fmt.Println("string_base64dec \"Some String to Encode\"")
            }
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "string_base64dec",
        Help: "Base64 Decode a string and print to screen.",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 1 {
                base64file.DecodeString(c.Args[0])
            } else {
                fmt.Println("string_base64dec <some string to decode>")
                fmt.Println(" EXAMPLE: >>> string_base64dec \"c3RyaW5nIHRvIGVuY29kZSBvciBzb21ldGhpbg==\"")
            }
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "zipdir",
        Help: "Zip up an entire directory recursively.",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 2 {
                zipdir.Zipdir(c.Args[0], c.Args[1])
            } else {
                fmt.Println("zipdir <zip file to create> <directory to zip>")
                fmt.Println(" EXAMPLE: >>> zipdir ./dirzip.zip ./dir")
            }
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "aes_encrypt_string",
        Help: "AES Encrypt a string and print it to the screen.",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 2 {
                fmt.Println(base64file.EncodeString(string(filecrypt.EncryptString([]byte(c.Args[1]), c.Args[0]))))
            } else {
                fmt.Println("aes_encrypt_string <passphrase> <string>")
            }
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "aes_decrypt_string",
        Help: "AES Decrypt a string and print it to the screen.",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 2 {
                fmt.Println(string(filecrypt.DecryptString([]byte(base64file.DecodeString(c.Args[1])), c.Args[0])))
            } else {
                fmt.Println("aes_decrypt_string <passphrase> <string>")
            }
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "pwd",
        Help: "Print working directory.",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 0 {
                mydir, err := os.Getwd()
                if err!= nil{fmt.Println(err)}
                fmt.Println("\t",mydir)
            } else {
                fmt.Println("pwd (no arguments)")
            }
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "nc_dash_l",
        Help: "Listen on a TCP port and hexdump everything received. (like 'nc -l', hence the name)",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 2 {
                port, err := strconv.Atoi(c.Args[1])
                if err==nil{
                    nc_dash_l.Nc_dash_l(c.Args[0], port)
                } else {
                    fmt.Println(err)
                }
            } else {
                fmt.Println("nc_dash_l <local_ip_to_bind_on> <port>")
            }
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "multicastchat",
        Help: "Start a multicast listener, for messages from other Seacrane LAN instances that multicastchat_send",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 1 {
                multicastchat.Multicastchat(c.Args[0])
            } else {
                fmt.Println("multicastchat <addr:port>")
                fmt.Println("\tExample: \n\tseacrane>> multicastchat 239.0.0.0:9999")
                fmt.Println("\nThe address in the example above works on most LANs btw.")
                fmt.Println("(Note: In reality multicast addr:port pairs are effectively just unique 'topics' that")
                fmt.Println("you are telling your router that you want to 'subscribe' to. Router must support it also.)\n")
            }
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "multicastchat",
        Help: "Start a multicast listener, for messages from other Seacrane LAN instances that multicastchat_send",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 1 {
                multicastchat.Multicastchat(c.Args[0])
            } else {
                fmt.Println("multicastchat <addr:port>")
                fmt.Println("\tExample: \n\tseacrane>> multicastchat 239.0.0.0:9999")
                fmt.Println("\nThe address in the example above works on most LANs btw.")
                fmt.Println("(Note: In reality multicast addr:port pairs are effectively just unique 'topics' that")
                fmt.Println("you are telling your router that you want to 'subscribe' to. Router must support it also.)\n")
            }
        },
    })

// QR code functions ONLY work on 64-bit platforms, not 32-bit like some of the mips targets.
// So commenting the next shell.AddCmd block will remove the qrcode submenu and code from teh tool.
// This needs to eventually be fixed with a compiletime directive or a git branch for mips-only releases.

// /* 
    shell.AddCmd(&ishell.Cmd{
        Name: "qrcode",
        Help: "Enter the QRCode submenu (has useful commands for encoding/decoding qrcodes).",
        Func: func(c *ishell.Context) {
            qrcode.Qrcodeshell()
        },
    })
// */

    shell.AddCmd(&ishell.Cmd{
        Name: "filedups",
        Help: "Recursive into directory finding files with identical contents. Also saves a 'directory fingerprint' log to disk.",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 1 {
                filedups.WalkDir(c.Args[0])
            } else {
                fmt.Println("filedups <directory>")
            }
        },
    })

    if len(os.Args) > 1 && os.Args[1] == "exit" {
        shell.Process(os.Args[2:]...)
    } else {
        // start shell
        shell.Run()
    }
}
