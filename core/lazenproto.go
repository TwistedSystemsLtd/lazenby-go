package core

import (
	"fmt"
	"github.com/TwistedSystemsLtd/lazenby-go/lazendata"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type LazenProto struct {}

func (p *LazenProto) CreateLazenFile(lazenfilePath string) {
	 createLazenProtoFile(lazenfilePath)
}

func (p *LazenProto)  SaveLazenFile(lazenfilePath string, lazenfile *lazendata.Lazenfile) {
	saveLazenProtoFile(lazenfilePath, lazenfile)
}

func (p *LazenProto) ReadLazenFile(lazenpath string) *lazendata.Lazenfile {
	return readLazenProtoFile(lazenpath)
}

func createLazenProtoFile(lazenfilePath string) {
	lazenkey := GenerateLazenkey()
	publicKey, privateKey := ReadUserKeys(Lazenhome())

	lazenkeys := make(map[string]*lazendata.Keypair)
	encryptedLazenKey := EncryptWithUserKey(publicKey, privateKey, lazenkey[:])

	DecryptWithUserKey(publicKey, privateKey, encryptedLazenKey)

	keypair := &lazendata.Keypair{PublicKey: publicKey[:], Lazenkey: encryptedLazenKey}
	lazenkeys[ToHexString(publicKey[:])] = keypair

	lazenfile := &lazendata.Lazenfile{Lazenkeys: lazenkeys, Secrets: nil}
	saveLazenProtoFile(lazenfilePath, lazenfile)
}


func saveLazenProtoFile(lazenfilePath string, lazenfile *lazendata.Lazenfile) {
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

func readLazenProtoFile(lazenpath string) *lazendata.Lazenfile {
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