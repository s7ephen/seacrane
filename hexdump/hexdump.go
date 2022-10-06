/* 
Hexdump

A tool wrapper around "hexout" which itself mostly a rip of 'hex_test' that was part of the `encoding/hex` package
found in the standard GoLang distribution.

This tool gets a filename from the user and then dumps its contens in a familiar 16-byte/per line "hexdump" style
output. 

*/

package hexdump

import (
//    "./hextest"
    "encoding/hex"
    "os"
    "fmt"
    "path/filepath"
    "io/ioutil"
//    "reflect"
)

var pln = fmt.Println
var pf = fmt.Printf
var p = fmt.Printf
var fpf = fmt.Fprintf

func Dump(file_name string) { 
//    fname := string(file_name)
//    fmt.Println("file_name:", file_name, reflect.TypeOf(file_name).Kind())
    fname := file_name
    abspath, err := filepath.Abs(fname)
    if err!=nil {
        pln("Unable to fetch absolute path of file. Error.")
        return 
    }
    if fname == "" {
        pln("Error: You must supply a file name!")
        return
    } 
    fdata, err := ioutil.ReadFile(fname)
/*
    // This is an attempt at testing the file size first, but it 
    file, errno := os.Open(*fname)
    defer file.Close()
    stat, err := file.Stat()
*/
    if err!=nil {
        pln(" [-] Error, unable to read file...",err)
        return
    }
    pln("\nBeginning the hexdump of ", abspath)
    stdoutDumper := hex.Dumper(os.Stdout)
    defer stdoutDumper.Close()
    stdoutDumper.Write(fdata)
}