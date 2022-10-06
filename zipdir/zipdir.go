package zipdir

import (
    "archive/zip"
    "fmt"
    "io"
    "os"
    "path/filepath"
)

func Zipdir(zipfilename string, dirtozip string) {
    file, err := os.Create(zipfilename)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    w := zip.NewWriter(file)
    fmt.Println("Creating zip file: ",zipfilename)
    defer w.Close()

    walker := func(path string, info os.FileInfo, err error) error {
        fmt.Printf("Crawling: %#v\n", path)
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
        file, err := os.Open(path)
        if err != nil {
            return err
        }
        defer file.Close()

        // Ensure that `path` is not absolute; it should not start with "/".
        // This snippet happens to work because I don't use 
        // absolute paths, but ensure your real-world code 
        // transforms path into a zip-root relative path.
        f, err := w.Create(path)
        if err != nil {
            return err
        }

        _, err = io.Copy(f, file)
        if err != nil {
            return err
        }

        return nil
    }
    err = filepath.Walk(dirtozip, walker)
    if err != nil {
        panic(err)
    }
}
