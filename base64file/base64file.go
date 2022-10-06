package base64file
/*

*/

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
    "bufio"
    "net/url"
    "path"
    "path/filepath"
)

func EncodeFile_nocompress(fname string) {
    // Open file on disk.
    f, err := os.Open(fname)
    if err != nil {
        fmt.Println(err)
    }
    defer f.Close()
    // Read entire JPG into byte slice.
    reader := bufio.NewReader(f)
    content, _ := ioutil.ReadAll(reader)
    // Encode as base64.
    encoded := base64.StdEncoding.EncodeToString(content)
    // Print encoded data to console.
    // ... The base64 image can be used as a data URI in a browser.
    fmt.Println("ENCODED: " + encoded)
    tempFile, err := ioutil.TempFile("./", fname+"-b64enc-")
    if err != nil {
        fmt.Println(err)
    }
    defer tempFile.Close()

    fmt.Println("\nWriting to: ", tempFile.Name())
    tempFile.Write([]byte(encoded))
}

func EncodeString(encode_string string) string{
    encoded := base64.StdEncoding.EncodeToString([]byte(encode_string))
//    fmt.Println("ENCODED: \n" + encoded)
    return encoded
}

func DecodeString(decode_string string) string{
    decoded, err := base64.StdEncoding.DecodeString(decode_string)
    if err != nil {
        fmt.Println(err)
    }
//    fmt.Println("DECODED: \n" + string(decoded))
    return string(decoded)
}

func URLEncode(encodestr string){
    fmt.Println("---URL ENCODED---\n", url.PathEscape(encodestr))
    fmt.Println("---QUERY ENCODED---\n", url.QueryEscape(encodestr))
}

func URLDecode(decodestr string){
    decoded, err := url.PathUnescape(decodestr)
    if err != nil {
        fmt.Println(err)
    }
    qdecoded, err := url.QueryUnescape(decodestr)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("---URL DECODED---\n", decoded)
    fmt.Println("---QUERY DECODED---\n", qdecoded)
}

func DecodeFile(fname string) {
    fname = filepath.Clean(fname) // app was segfaulting if user provided stuff like ./
	var text []byte
	fileContent, err := os.ReadFile(fname)
	text = fileContent

	result, err := decode(text)
    if err != nil {
        fmt.Println("\t[+] Decode failed...")
        return
    }

	fmt.Print("success!\n\n")
    curdir, err := os.Getwd()
    tempFile, err := ioutil.TempFile(curdir, path.Base(fname)+"-b64dec-")
    if err != nil { fmt.Println(err.Error()) }
    defer tempFile.Close()

    fmt.Println("\nWriting to: ", tempFile.Name())
    tempFile.Write([]byte(result))

	// output it
	//os.Stdout.Write([]byte(result))
/*
	fmt.Print("\n\n")
    fmt.Println("Outputting to: ",fname+".b64decoded")
	os.WriteFile(fname+".b64decoded", []byte(result), os.FileMode(777))
*/
}

func EncodeFile(fname string) {
    fname = filepath.Clean(fname) // app was segfaulting if user provided stuff like ./
//	fileContent, err := os.ReadFile(fname)
    inFile, err := os.Open(fname)
    if err != nil { fmt.Println(err.Error()) }
    defer inFile.Close()
    reader := bufio.NewReader(inFile)
    content, _ := ioutil.ReadAll(reader)

	result, err := encode(content)
	if err != nil {
		fmt.Println("\t[+] Encode Failed...")
		return
	}
    curdir, err := os.Getwd()
    tempFile, err := ioutil.TempFile(curdir,path.Base(fname)+"-b64enc-")
    if err != nil { fmt.Println(err.Error()) }
    defer tempFile.Close()

	// output it
//	os.Stdout.Write([]byte(result))
    fmt.Println("\nWriting to: ", tempFile.Name())
    tempFile.Write([]byte(result))
/*
	fmt.Print("\n\n")
    fmt.Println("Outputting to: ",fname+".b64encoded")
	os.WriteFile(fname+".b64encoded", []byte(result), os.FileMode(777))
*/
}

func decode(text []byte) (result string, err error) {
//Compresses and encodes.

	fmt.Print("Attempting to decode the content... ")

	// base64 decode it
	textDecoded := make([]byte, len(text))
	_, err = base64.RawStdEncoding.Decode(textDecoded, text)
	if err != nil {
		fmt.Print(err)
		return
	}

	// decompress it
	reader := bytes.NewReader(textDecoded)
	gzreader, err := gzip.NewReader(reader)
	gzreader.Multistream(false)
	if err != nil {
		fmt.Print("error at stage 1: " + err.Error())
		return
	}

	resultBytes, err := ioutil.ReadAll(gzreader)
	if err != nil {
		fmt.Print("error at stage 2: " + err.Error())
	}

	result = string(resultBytes)
	return
}

func encode(text []byte) (result string, err error) {
// Decompresses and Decodes.

	fmt.Print("\nAttempting to encode the content... ")

	// compress it
	buf := new(bytes.Buffer)
	gz := gzip.NewWriter(buf)
	_, err = gz.Write(text)
	if err != nil {
		fmt.Print(err)
		return
	}
	gz.Close()

	// base64 encode it
	result = base64.RawStdEncoding.EncodeToString(buf.Bytes())

	return
}

func main(){
    EncodeFile(os.Args[1])
//    DecodeFile(os.Args[1])
}
