package core

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/TwistedSystemsLtd/lazenby-go/lazendata"
	"io/ioutil"
	"log"
	"os"
)

type LazenToml struct {}

func (p *LazenToml) CreateLazenFile(lazenfilePath string) {
	createLazenTomlFile(lazenfilePath)
}

func (p *LazenToml)  SaveLazenfile(lazenfilePath string, lazenfile *lazendata.Lazenfile) {
	saveLazenTomlFile(lazenfilePath, lazenfile)
}

func (p *LazenToml) ReadLazenFile(lazenpath string) *lazendata.Lazenfile {
	return readLazenTomlFile(lazenpath)
}


func createLazenTomlFile(lazenfilePath string) {
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
	saveLazenTomlFile(lazenfilePath, lazenfile)
}

func saveLazenTomlFile(lazenfilePath string, lazenfile *lazendata.Lazenfile) {
	buf := new(bytes.Buffer)
	marshaler := toml.NewEncoder(buf)
	err := marshaler.Encode(lazenfile)
	if err != nil {
		log.Panic("Error marshalling lazenfile", err)
	}
	writeErr := ioutil.WriteFile(lazenfilePath, buf.Bytes(), 0644)
	if writeErr != nil {
		log.Panic("Could not write lazenfile", writeErr)
	}
}

func readLazenTomlFile(lazenpath string) *lazendata.Lazenfile {
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
	toml.Unmarshal(lazenBodyData, parsedLazenfile)

	return parsedLazenfile
}