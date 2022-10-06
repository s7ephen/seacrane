/*
Finds duplicate files in a directory tree by walking the tree recursively.
Also sorts the results and reports them as well as writing a "directory fingerprint"
log of all files and their hashes to disk.
*/
package filedups

import "fmt"
import "time"
import "os"
import "io"
import "path/filepath"
import "crypto/sha256"
import "sort"
import "encoding/json"
import "io/ioutil"

type filehashinfo struct {
    Path   string
    Hash   string
    Size    int64
    TimeModified time.Time
}

func main(){
    WalkDir("./")
}

func getSHA(path string) string {
    f, err := os.Open(path)
    if err != nil {
        fmt.Println(err)
        return ""
    }

    defer f.Close()
    h := sha256.New()
    if _, err := io.Copy(h, f); err != nil {
        fmt.Println("--",err,"--")
    }
//    fmt.Println("%x", h.Sum(nil))
    return string(fmt.Sprintf("%x",h.Sum(nil)))
}

func WalkDir(targetpath string) {
    filehashlist := []string{}
    filetree := map[string]filehashinfo{}
    occurence_count := map[string]int{}
    err := filepath.Walk(targetpath,
        func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        //fmt.Println(path, info.Size())
        if info.Mode().IsRegular() { // check if it is a regular file
            hash := getSHA(path)
            fhi := filehashinfo{Path: path, Hash: hash, Size: info.Size(), TimeModified: info.ModTime()}
            filehashlist = append(filehashlist, hash)
            filetree[path]=fhi
/*
            fmt.Println("---")
            fmt.Println("\tFile: ",path, "\n\tFile Size: ", info.Size(), "\n\tHash: ",hash)
//            fmt.Printf("%+v", fhi) //%+v prints the elements of the struct AND their value!
//            fmt.Printf("%#v", fhi) //prints the elements of the struct, their values, AND the full name of the struct
*/
        }
        return nil
    })
    if err != nil {
        fmt.Println(err)
    }
    for _, fhash := range filehashlist {
        occurence_count[fhash]=occurence_count[fhash]+1
    }
// ---- 
    jsonStr, err := json.MarshalIndent(filetree,"","\t")
//    jsonStr, err := json.Marshal(filetree)
    if err != nil {
        fmt.Printf("Error: %s", err.Error())
    } /*else {
        fmt.Println(string(jsonStr))
    }*/

/*
    fmt.Println("\n----filetree---\n",filetree)
    fmt.Println("\n----Occurrences---\n",occurence_count)
*/

// Bear with me, this is ugly...to sort the occurence_count map by value:
// Since golang doesnt have the same kinds of .has_element() or .has_key() convenience functions
// that Python has we have to do some quasi-C shit here with the map, putting it in a separate map
// keyed on values so we can sort....we also have do do stuff accounting for our typedef struct
    occurence_count_keys := make([]string, 0, len(occurence_count))
    for k := range occurence_count{
        occurence_count_keys = append(occurence_count_keys, k)
    }
    sort.SliceStable(occurence_count_keys, func(i, j int) bool{
        //return occurence_count[occurence_count_keys[i]] > occurence_count[occurence_count_keys[j]]
        //flip the > to a < above for sorting from least to greatest.
        return occurence_count[occurence_count_keys[i]] < occurence_count[occurence_count_keys[j]]
        //flip the < to a > above for sorting from greatest to least.
    })
//    fmt.Println("\n----Occurences SORTED---\n")
    for _, k := range occurence_count_keys {
        //fmt.Println(k, occurence_count[k])
        prev:="" //since we are walking through a sorted list we can use this to stash the previous
                //iteration's value to just append to it instead of printing out a whole new block of file duplicate 
                //info
        if occurence_count[k] > 1{
            //SOOO UGLY, now we need to find this hash in the tree to give a sample
            //filename to the user to show its contents are repeeated inside other files.
            //This shoudl be a helper function but...
            for gethashfname := range filetree { //walk through filetree looking at fhi structs
                if filetree[gethashfname].Hash == prev {
                    fmt.Println("\tDUPLICATE FILE: ", filetree[gethashfname].Path)
                } else if filetree[gethashfname].Hash == k {
                    fmt.Println("-----")
                    fmt.Println("The contents of a file has duplicates!")
                    fmt.Println("\tFILE: ", filetree[gethashfname].Path)
                    fmt.Println("\tHASH: ", filetree[gethashfname].Hash)
                    fmt.Println("\tDUPLICATE COUNT: ",occurence_count[k])
                    prev = k
                }
            }
        }
    }
// ---- end "bear with me"----
    curdir, err := os.Getwd()
    tempFile, err := ioutil.TempFile(curdir,"directory_fingerprint.json-")
    if err != nil { fmt.Println(err.Error()) }
    defer tempFile.Close()
    fmt.Println("\nWriting complete directory fingerprint to: ", tempFile.Name())
    tempFile.Write(jsonStr)
}
