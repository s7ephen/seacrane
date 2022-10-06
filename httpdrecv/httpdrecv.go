package httpdrecv

import (
    "fmt"
    "io/ioutil"
    "net/http"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
    fmt.Println("File Upload Endpoint Hit")

    // Parse our multipart form, 10 << 20 specifies a maximum
    // upload of 10 MB files.
    r.ParseMultipartForm(10 << 20)
    // FormFile returns the first file for the given key `myFile`
    // it also returns the FileHeader so we can get the Filename,
    // the Header and the size of the file
    file, handler, err := r.FormFile("myFile")
    if err != nil {
        fmt.Println("Error Retrieving the File")
        fmt.Println(err)
        return
    }
    defer file.Close()
    fmt.Printf("Uploaded File: %+v\n", handler.Filename)
    fmt.Printf("File Size: %+v\n", handler.Size)
    fmt.Printf("MIME Header: %+v\n", handler.Header)

    // Create a temporary file within our temp-images directory that follows
    // a particular naming pattern
    tempFile, err := ioutil.TempFile("./", handler.Filename+"-")
    if err != nil {
        fmt.Println(err)
    }
    defer tempFile.Close()

    // read all of the contents of our uploaded file into a
    // byte array
    fileBytes, err := ioutil.ReadAll(file)
    if err != nil {
        fmt.Println(err)
    }
    // write this byte array to our temporary file
    tempFile.Write(fileBytes)
    fmt.Println("\t[+] Received file and written to: ", tempFile.Name())
    // return that we have successfully uploaded our file!
    fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func printIndex(w http.ResponseWriter, r *http.Request) {
    indexhtml :=` 
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta http-equiv="X-UA-Compatible" content="ie=edge" />
    <title>Document</title>
  </head>
  <body>
    <form
      enctype="multipart/form-data"
      action="/upload"
      method="post"
    >
      <input type="file" name="myFile" />
      <input type="submit" value="upload" />
    </form>
  </body>
</html>
`
    fmt.Fprintf(w, indexhtml)
}

func StartServer(ip string, port string) {
    http.HandleFunc("/upload", uploadFile)
    http.HandleFunc("/", printIndex)
    fmt.Println("\t[+] Started HTTP File upload server on: ", ip+":"+port)
    http.ListenAndServe(ip+":"+port, nil)
}

func main() {
    StartServer("0.0.0.0","8080")
}
