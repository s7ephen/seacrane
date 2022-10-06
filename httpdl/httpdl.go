package httpdl

import (
	"fmt"
    "time"
	"io"
	"net/http"
	"os"
	"strings"
    "github.com/dustin/go-humanize" //<--- just something for printing out mb/Mb/Gig etc in human-readable.
)

type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
// An overload for the Write method. uses the WriteCounter struct to store
// a byte counter (for number of bytes written) which will report progress via
// io.TeeReader() with each write invocation
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\r\tDownloading... %s complete", humanize.Bytes(wc.Total))
}

func DoDownload(filename string, fileUrl string) {
	fmt.Println("\t[+] Download Started at: "+time.Now().String())
	err := DownloadFile(filename, fileUrl)
	if err != nil {
		panic(err)
	}
	fmt.Println("\t[+] Download Finished at: "+time.Now().String())
}

// Download a file to disk directly without buffering to memory.
// do this by passing an io.TeeReader into Copy() to report progress on the download.
func DownloadFile(filepath string, url string) error {

	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()

	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}
	fmt.Print("\n")
	out.Close()

	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}
	return nil
}
