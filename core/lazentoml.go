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
	lazenfile := NewLazenFile()
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