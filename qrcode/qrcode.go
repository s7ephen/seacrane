package qrcode

import (
	"log"
    "os"
    "github.com/abiosoft/ishell"
    "github.com/YuriyLisovskiy/qrcode/qr"
    "github.com/tuotoo/qrcode"
)


func Qrcodeshell() {
    shell := ishell.New()
    shell.SetPrompt("seacrane>> qrcode>> ")
    // display welcome info.
    shell.Println("QRCode submenu")
    shell.AddCmd(&ishell.Cmd{
        Name: "qrdecode_file",
        Help: "Decodes a qrcode from a png file and prints the decoder results to the screen. ",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 1 {
                Qrcode_decode_file(c.Args[0])
            } else { c.Println("qrdecode_file <file_to_decode>")}
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "qrencode_to_file",
        Help: "Encodes a string to a QRCode and stores it to disk as PNG data.",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 1 {
                c.ShowPrompt(false)
                defer c.ShowPrompt(true) // yes, revert after input 
                c.Println("Input however many lines you want, just end with ';;;'.")
                lines := c.ReadMultiLines(";;;")
                c.Println("\t[+] Got it.")
                Qrcode_encode_to_file(c.Args[0],lines)
            } else { c.Println("qrencode_to_file <file_to_save_to> (you are prompted for the data to encode)")}
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "qrencode_print",
        Help: "Encodes a string to a QRCode and prints the raw PNG data to screen (some terminals actually render it).",
        Func: func(c *ishell.Context) {
            if len(c.Args) >= 0 {
                c.ShowPrompt(false)
                defer c.ShowPrompt(true) // yes, revert after input 
                c.Println("Input however many lines you want, just end with ';;;'.")

                lines := c.ReadMultiLines(";;;")
                c.Println("\t[+] Got it.")
                Qrcode_print(lines)
            } else { c.Println("qrencode_print (no arguments, you are prompted for the data to encode)")}
        },
    })
    shell.AddCmd(&ishell.Cmd{
        Name: "file2qrcode",
        Help: "INCOMPLETE: Reads a file from disk, compresses, and stores it as a qrcode if it is small enough.",
        Func: func(c *ishell.Context) {
            c.Println("\nUnder construction")
        },
    })
    shell.Run()
}

func Qrcode_print(TEXT string) {
	qrGenerator := qr.Generator{}
	qrGenerator = qrGenerator.EncodeText(TEXT)
//	qrGenerator.DrawImage("qr.png", 4, 500)
	qrGenerator.Draw(4)
}

func Qrcode_encode_to_file(filename string, TEXT string) {
    log.Println("\t[+] Attempting to encode and write to : ", filename)
	qrGenerator := qr.Generator{}
	qrGenerator = qrGenerator.EncodeText(TEXT)
	qrGenerator.DrawImage(filename, 4, 500)
//	qrGenerator.Draw(4)
}

func Qrcode_decode_file(filename string){
    log.Println("\t[+] Attempting to decode from : ", filename)
    fi, err := os.Open(filename)
    if err != nil{
        log.Println(err.Error())
        return
    }
    defer fi.Close()
    qrmatrix, err := qrcode.Decode(fi)
    if err != nil{
        log.Println(err.Error())
        return
    }
    log.Println(qrmatrix.Content)
}
