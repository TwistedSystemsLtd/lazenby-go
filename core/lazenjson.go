package core

import (
	"bytes"
	"fmt"
	"github.com/TwistedSystemsLtd/lazenby-go/lazendata"
	"github.com/golang/protobuf/jsonpb"
	"io/ioutil"
	"log"
	"os"
)

type LazenJson struct {}

func (p *LazenJson) CreateLazenFile(lazenfilePath string) {
	createLazenJsonFile(lazenfilePath)
}

func (p *LazenJson)  SaveLazenFile(lazenfilePath string, lazenfile *lazendata.Lazenfile) {
	saveLazenJsonFile(lazenfilePath, lazenfile)
}

func (p *LazenJson) ReadLazenFile(lazenpath string) *lazendata.Lazenfile {
	return readLazenJsonFile(lazenpath)
}


func createLazenJsonFile(lazenfilePath string) {
	lazenkey := GenerateLazenkey()
	publicKey, privateKey := ReadUserKeys(Lazenhome())

	lazenkeys := make(map[string]*lazendata.Keypair)
	encryptedLazenKey := EncryptWithUserKey(publicKey, privateKey, lazenkey[:])

	DecryptWithUserKey(publicKey, privateKey, encryptedLazenKey)

	currentUser := CurrentUser()
	username := currentUser.Username
	name := currentUser.Name

	tags := []string{username}
	if  name  != "" {
		tags = append(tags, name)
	}

	keypair := &lazendata.Keypair{PublicKey: publicKey[:], Lazenkey: encryptedLazenKey, Tags: tags}
	lazenkeys[ToHexString(publicKey[:])] = keypair

	lazenfile := &lazendata.Lazenfile{Lazenkeys: lazenkeys, Secrets: nil}
	saveLazenJsonFile(lazenfilePath, lazenfile)
}

func saveLazenJsonFile(lazenfilePath string, lazenfile *lazendata.Lazenfile) {
	marshaler := jsonpb.Marshaler{Indent: "  "}
	lazenbytes, err := marshaler.MarshalToString(lazenfile)
	if err != nil {
		log.Panic("Error marshalling lazenfile", err)
	}
	writeErr := ioutil.WriteFile(lazenfilePath, []byte(lazenbytes), 0644)
	if writeErr != nil {
		log.Panic("Could not write lazenfile", writeErr)
	}
}

func readLazenJsonFile(lazenpath string) *lazendata.Lazenfile {
	if _, err := os.Stat(lazenpath); err == nil {
		log.Print(fmt.Sprintf("Lazenfile exists"))
	} else {
		log.Panic("Lazenfile not found", lazenpath)
	}
	lazenBodyData, readErr := ioutil.ReadFile(lazenpath)

	if readErr != nil {
		log.Panic("Error reading lazenpath", readErr)
	}

	parsedLazenfile := &lazendata.Lazenfile{}
	jsonpb.Unmarshal(bytes.NewReader(lazenBodyData), parsedLazenfile)

	return parsedLazenfile
}