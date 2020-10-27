package main

import (
	"bytes"
	"crypto/rc4"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/folbricht/pefile"
	. "github.com/logrusorgru/aurora"
)

// check errors as they occur and panic :o
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// rc4decrypt acts as a helper function for RC4-XOR decryption
func rc4decrypt(extkey []byte, data []byte) []byte {
	// create a new RC4 Enc/Dec Routine and pass the key
	cipher, ciphErr := rc4.NewCipher(extkey)
	check(ciphErr)
	// decrypt the config
	cipher.XORKeyStream(data, data)

	return data
}

func main() {

	fmt.Printf("\n\n     ,-----.,------. ,----.   ,--.   ,--.        ,--.,--.                  \n")
	fmt.Printf("    '  .--./|  .---''  .-./   |  |   |  | ,--,--.|  ||  |,-. ,---. ,--.--. \n")
	fmt.Printf("    |  |    |  `--, |  | .---.|  |.'.|  |' ,-.  ||  ||     /| .-. :|  .--' \n")
	fmt.Printf("    '  '--'\\|  |`   '  '--'  ||   ,'.   |\\ '-'  ||  ||  \\  \\    --.|  |\n")
	fmt.Printf("     `-----'`--'     `------' '--'   '--' `--`--'`--'`--'`--'`----'`--' \n")
	fmt.Printf("\n                Netwalker Ransomware Configuration Extractor\n")
	fmt.Printf("            Marius 'f0wL' Genheimer | https://dissectingmalwa.re\n\n\n")

	if len(os.Args) < 2 {
		fmt.Println("   Usage: go run cfgwalker.go netwalker_sample.exe [--print (write to stdout instead of a file)]")
		os.Exit(1)
	}

	// open the sample with pefile
	f, openErr := pefile.Open(os.Args[1])
	check(openErr)
	defer f.Close()
	// Fetch resources from the PE
	resources, resErr := f.GetResources()
	check(resErr)

	// Length of the RC4 key
	var keyLen int
	// RC4 Key needed for decryption of the resource
	var rc4Key []byte
	// holds the encrypted resource
	var encryptedData []byte

	for _, r := range resources {
		// The Netwalker config is safed in the Resource section under 1337/31337
		if r.Name == "1337/31337/0" {
			encryptedData = r.Data
		} else {
			fmt.Printf(Sprintf(Red(" [!] Error: No config resource found.\n\n")))
		}
	}

	if len(os.Args) < 2 {
		// write encrypted resource to file for safe-keeping
		writeErr := ioutil.WriteFile("config.enc", encryptedData, 0644)
		check(writeErr)
		fmt.Printf(Sprintf(Green(" [+] Wrote extracted resource to %s"), White("'config.enc'\n\n")))
	}

	// the length of the RC4 Key is defined in the first 4 bytes of the file
	keyLen = int(encryptedData[0])
	fmt.Printf(Sprintf(Green(" [+] Key Length: ")))
	fmt.Printf("%v\n\n", keyLen)

	// we'll skip the first 4 bytes and retrieve the encryption key using the previously determined length
	rc4Key = encryptedData[4 : 4+keyLen]
	fmt.Printf(Sprintf(Green(" [+] Extracted RC4 Key: ")))
	fmt.Printf("%v\n\n", hex.EncodeToString(rc4Key))

	// decrypt the RC4 encrypted config
	config := rc4decrypt(rc4Key, encryptedData[4+keyLen:])

	// beautify the json string
	var out bytes.Buffer
	jsonErr := json.Indent(&out, config, "", "  ")
	check(jsonErr)

	if len(os.Args) > 2 && os.Args[2] == "--print" {
		fmt.Printf(Sprintf(Green(" [+] Decrypted JSON Config: \n\n")))
		fmt.Printf("%v\n\n", out.String())
	} else {
		// Write json config to a file
		writeErr := ioutil.WriteFile("config.json", out.Bytes(), 0644)
		check(writeErr)
		fmt.Printf(Sprintf(Green(" [+] Wrote decrypted config to %s"), White("'config.json'\n\n")))
	}

	// and now for some fancy Golang magic :D Extracting the Ransomnote from the json string via an Interface
	// This saves a lot of time and is less prone to errors since we don't have to unmarshal into a structure
	// Declaring an empty interface
	var confInterf map[string]interface{}

	// unmarshal the json string into the interface
	json.Unmarshal([]byte(out.String()), &confInterf)

	// extract the encoded ransomnote string
	ransomnote := confInterf["lend"]
	// decode the note from base64, var ransomnote is type-asserted to string
	decodedNote, base64Err := base64.StdEncoding.DecodeString(ransomnote.(string))
	check(base64Err)

	if len(os.Args) > 2 && os.Args[2] == "--print" {
		fmt.Printf(Sprintf(Green("\n\n [+] Extracted and base64 decoded Ransomnote:\n\n")))
		fmt.Printf(string(decodedNote) + "\n\n")
	} else {
		// Write decoded ransomnote text to a file
		writeErr := ioutil.WriteFile("ransomnote.txt", decodedNote, 0644)
		check(writeErr)
		fmt.Printf(Sprintf(Green(" [+] Wrote decoded ransomnote to %s"), (White("'ransomnote.txt'\n\n"))))
	}
}
