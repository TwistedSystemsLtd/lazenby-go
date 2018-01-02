package core

import (
	"github.com/TwistedSystemsLtd/lazenby-go/lazendata"
	"log"
	"os"
	"path"
)

type LazenMarshaller interface {
	CreateLazenFile(lazenfilePath string)
	SaveLazenFile(lazenfilePath string, lazenfile *lazendata.Lazenfile)
	ReadLazenFile(lazenfilePath string) *lazendata.Lazenfile
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

func CreateLazenFile(lazenfilePath string) {
	lf := new(LazenToml)
	lf.CreateLazenFile(lazenfilePath)
}
func SaveLazenFile(lazenfilePath string, lazenfile *lazendata.Lazenfile) {
	lf := new(LazenToml)
	lf.SaveLazenfile(lazenfilePath, lazenfile)
}
func ReadLazenFile(lazenfilePath string) *lazendata.Lazenfile {
	lf := new(LazenToml)
	return lf.ReadLazenFile(lazenfilePath)
}

func NewLazenFile() *lazendata.Lazenfile {
	lazenkey := GenerateLazenkey()
	publicKey, privateKey := ReadUserKeys(Lazenhome())

	lazenkeys := make(map[string]string)
	encryptedLazenKey := EncryptWithUserKey(publicKey, privateKey, lazenkey[:])

	DecryptWithUserKey(publicKey, privateKey, encryptedLazenKey)

	currentUser := CurrentUser()
	username := currentUser.Username
	name := currentUser.Name

	tags := []string{username}
	if  name  != "" {
		tags = append(tags, name)
	}

	keypair := EncodeString(encryptedLazenKey)
	lazenkeys[EncodeString(publicKey[:])] = keypair

	return &lazendata.Lazenfile{Lazenkeys: lazenkeys, Secrets: make(map[string]string)}
}