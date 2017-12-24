package core

import (
	"log"
	"io/ioutil"
	"github.com/TwistedSystemsLtd/lazenby-go/lazendata"
	"github.com/golang/protobuf/proto"
	"strings"
	"os"
	"fmt"
	"path"
)

func CreateLazenfile(lazenfilePath string) {
	lazenkey := GenerateLazenkey()
	publicKey, privateKey := ReadUserKeys(Lazenhome())

	lazenkeys := make(map[string]*lazendata.Keypair)
	encryptedLazenKey := EncryptWithUserKey(publicKey, privateKey, lazenkey[:])

	DecryptWithUserKey(publicKey, privateKey, encryptedLazenKey)

	keypair := &lazendata.Keypair{PublicKey: publicKey[:], Lazenkey: encryptedLazenKey}
	lazenkeys[ToHexString(publicKey[:])] = keypair

	lazenfile := &lazendata.Lazenfile{Lazenkeys: lazenkeys, Secrets: nil}
	SaveLazenfile(lazenfilePath, lazenfile)
}

func SaveLazenfile(lazenfilePath string, lazenfile *lazendata.Lazenfile) {
	lazenbytes, err := proto.Marshal(lazenfile)
	if err != nil {
		log.Panic("Error marshalling lazenfile", err)
	}
	hexString := ToHexString(lazenbytes)
	chunks := Chunk(hexString)
	body := []byte(strings.Join(chunks, "\n"))
	writeErr := ioutil.WriteFile(lazenfilePath, body, 0644)
	if writeErr != nil {
		log.Panic("Could not write lazenfile", writeErr)
	}
}

func ReadLazenfile(lazenpath string) *lazendata.Lazenfile {
	if _, err := os.Stat(lazenpath); err == nil {
		log.Print(fmt.Sprintf("Lazenfile exists"))
	} else {
		log.Panic("Lazenfile not found", lazenpath)
	}
	lazenBodyData, readErr := ioutil.ReadFile(lazenpath)

	if readErr != nil {
		log.Panic("Error reading lazenpath", readErr)
	}

	lazenhex := strings.Replace(string(lazenBodyData), "\n", "", -1)

	lazenbytes := FromHexString(lazenhex)

	parsedLazenfile := &lazendata.Lazenfile{}
	proto.Unmarshal(lazenbytes, parsedLazenfile)

	return parsedLazenfile
}

func GetLazenpath(lazenfile string) string {
	var lazenpath string
	dir, err := os.Getwd()
	if err != nil {
		log.Panic("Could not get current working directory")
	}
	if path.IsAbs(lazenfile) {
		lazenpath = lazenfile
	} else {
		lazenpath = path.Join(dir, lazenfile)
	}
	return lazenpath
}

