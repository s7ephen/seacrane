package findurls
import (
    "fmt"
    "regexp"
    "io/ioutil"
    "os"
)

func CheckError(e error) {
    if e != nil {
        fmt.Println(e)
    }
}

func FindURLs(infilename string, outfile string) {
    infile, err := ioutil.ReadFile(infilename)
    fmt.Println("\t[+] Printing URLs found in `%s` to screen and saving to '%s'", infilename, outfile)
	if err != nil {
		fmt.Printf("\n[-]\t%v\n", err)
	}

    f, err := os.Create(outfile)
    if err != nil {
    fmt.Println(err)
    }
    defer f.Close()

    re, e := regexp.Compile(`(?:(?:https?|ftp):\/\/)?[\w/\-?=%.]+\.[\w/\-&?=%.]+`)
    CheckError(e)
    print_count := 0
    for _, value:= range re.FindAll([]byte(infile),-1){
        if print_count <= 10 {
            fmt.Println(string(value))
            if print_count == 10 { fmt.Println("\n[truncating match output]\n") }
        } else {
            fmt.Print(".")
        }
        f.Write([]byte(string(value)))
        print_count++
    }
    fmt.Println("\n\t[+] Done printing URLs found in `%s` to screen and saving to '%s'", infilename, outfile)
    if print_count >= 10 {
        fmt.Println("\t\t[+] Screen output truncated ... but `%d` matches saved to `%s`", print_count, outfile)
    }
}

func main(){
    FindURLs("./test.html", "./foundurls.txt")
}

