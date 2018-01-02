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
