package filecrypt 
/*
This will symmetrically encrypt a file (password.)

There is a lot of extra stuff in here to be used later:
    -encryptString() / decryptString() was when I was playing with crypto functions
    -encryptToFileTest() / decryptFromFileTest() was when I was trying to do file crypto, but it reads into memory first
    -someTests() just calls the above to play with them.
    -encryptFromFile() was me playing with File streams, but it still buffers into memory, I also didnt bother with
        a decryption function after I realized that io.Copy() was still buffering in memory the way I did it.
    -everything else below the above is functional Stream-based encryption/decryption

TODOS:
  See "TODO/NOTES:: in the section below.

*/
import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"bytes"
    "strings" // <---- only for my lazy file string handling.
)

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func EncryptString(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func DecryptString(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func encryptToFileTest(filename string, data []byte, passphrase string) {
	f, _ := os.Create(filename)
	defer f.Close()
	f.Write(EncryptString(data, passphrase))
}

func decryptFromFileTest(filename string, passphrase string) []byte {
	data, _ := ioutil.ReadFile(filename)
	return DecryptString(data, passphrase)
}

/*
Run all of the above to test out string based encryption. 
Everything below the separator below is for file stream based 
crypto, and all the above will be integrated more cleanly later.
*/
func someTests() {
	fmt.Println("Starting the application...")
	ciphertext := EncryptString([]byte("Hello World"), "password")
	fmt.Printf("Encrypted: %x\n", ciphertext)
	plaintext := DecryptString(ciphertext, "password")
	fmt.Printf("Decrypted: %s\n", plaintext)
	encryptToFileTest("sample.txt", []byte("Hello World"), "password1")
	fmt.Println(string(decryptFromFileTest("sample.txt", "password1")))
}

// ---------------Everything Below is for file encrpytion-------------------

func main() {
	//encryptFromFile()

//   --- These below are test cases. Still too lazy to learn the built-in Golang unit tests --- 
//    StreamEncrypter("passphrase","InstagramForChimps.gif")
//    StreamDecrypter("passphrase","InstagramForChimps.gif.sa7")
}

/*
Even though it works, the problem with this test encrpytion routine is that
it uses File streaming but doesnt do it without  reading full contents into memory. So it sucks for large files.
I also didnt write a decrpytion routine before I realize this. Saving for later...
*/
func encryptFromFile() {
    // read content from your file
	fmt.Println("Starting encryption run...")
    plaintext, err := ioutil.ReadFile("plaintextfile")
    if err != nil {
        panic(err.Error())
    }

    key := []byte("this_is_the_pass_word")

    block, err := aes.NewCipher(key)
    if err != nil {
        panic(err.Error())
    }

    // prepend a unique IV (which doesnt need to be secret)
    ciphertext := make([]byte, aes.BlockSize+len(plaintext))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        panic(err.Error())
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

    // create a new file for saving the encrypted data.
    f, err := os.Create("encryptedfile")
    if err != nil {
        panic(err.Error())
    }
    _, err = io.Copy(f, bytes.NewReader(ciphertext))
    if err != nil {
        panic(err.Error())
    }
}

/* ----------------------- EVERYTHING below is function stream-based file encryption/decryption -----------------
TODO/NOTES:
 
- IV SHIT: 
    This method uses are zero'd IV (unlike the in-memory Read method in encryptFromFile() ) mostly because I am lazy, 
    and didnt want to embed and parse out the IV from seek()ing into the file without reading it into memory.
    THEREFORE, this needs to be updated to support embedding the IV. 
    For this current nulled IV to be secure the key must be unique for each ciphertext to prevent a Bruteforce on the 
    ciphertext

 - AUTHENTICATION / CHECKSUMMING
    This lazy method also does not authenticate the encrypted data. Attacker *could* attack ciphertext by doing
    attacks on the ciphertext (random bitflipping) to leak/disclose the key eventually if it was for something
    continuous like a TCP stream, but for files this should be fine, especially if the key is different for each
    ciphertext.
*/


func StreamDecrypter(passphrase string, file string) {
//  key must be multiples of 8 bytes to initialize the NewCipher. to achieve this we use createHash()
//  so the passphrase used for encyption is not the verbatim passphrase, it is instead a md5 of it. 
    key := []byte(createHash(passphrase))
    fmt.Println(len(passphrase), "---", len(createHash(passphrase)))

    inFile, err := os.Open(file)
    if err != nil {
        panic(err.Error())
//        panic(err)
    }
    defer inFile.Close()

    // NewCipher aes requires key aligned on four-word boundary. 
    block, err := aes.NewCipher(key)
    if err != nil {
        panic(err.Error())
//        panic(err)
    }

    // NOTE THE NULLED IV below! Ciphertext well therefore need unique key to avoid BF attacks on the ciphertext.
    var iv [aes.BlockSize]byte
    stream := cipher.NewOFB(block, iv[:])

    outFile, err := os.OpenFile(strings.Split(file,".sa7")[0]+".dec", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
    fmt.Println("\t[+] Outputting to: ",strings.Split(file,".sa7")[0]+".dec")
    if err != nil {
        panic(err.Error())
//        panic(err)
    }
    defer outFile.Close()

//  The Golang Stream magic is on the line below. 
    reader := &cipher.StreamReader{S: stream, R: inFile}
    // Do the copy decrypting as we go without slurping into memory.
    if _, err := io.Copy(outFile, reader); err != nil {
        panic(err.Error())
//        panic(err)
    }
}

func StreamEncrypter(passphrase string, file string) {
//  key must be multiples of 8 bytes to initialize the NewCipher. to achieve this we use createHash()
//  so the passphrase used for encyption is not the verbatim passphrase, it is instead a md5 of it. 
    key := []byte(createHash(passphrase))
    fmt.Println(len(passphrase), "---", len(createHash(passphrase)))
    inFile, err := os.Open(file)
    if err != nil {
        panic(err.Error())
//        panic(err)
    }
    defer inFile.Close()

    // NewCipher aes requires key aligned on four-word boundary. 
    block, err := aes.NewCipher(key)
    if err != nil {
        panic(err.Error())
//        panic(err)
    }

    // NOTE THE NULLED IV below! Ciphertext well therefore need unique key to avoid BF attacks on the ciphertext.
    var iv [aes.BlockSize]byte
    stream := cipher.NewOFB(block, iv[:])

//    outFile, err := os.OpenFile("encrypted-file", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
    outFile, err := os.OpenFile(file+".sa7", os.O_WRONLY|os.O_CREATE, 0777) // fuck it make it world readable or something.
    fmt.Println("\t[+] Outputting to: ",file+".sa7")
    if err != nil {
        panic(err.Error())
//        panic(err)
    }
    defer outFile.Close()

//  The Golang Stream magic is on the line below. 
    writer := &cipher.StreamWriter{S: stream, W: outFile}
    // Do the copy encrypting as we go without slurping into memory.
    if _, err := io.Copy(writer, inFile); err != nil {
        panic(err.Error())
//        panic(err)
    }
}
