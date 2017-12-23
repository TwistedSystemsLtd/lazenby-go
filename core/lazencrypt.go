package core

import (
	"io"
	"golang.org/x/crypto/nacl/secretbox"
	"crypto/rand"
	"log"
	"fmt"
	"encoding/hex"
	"golang.org/x/crypto/nacl/box"
)

func genKey() [32]byte {
	var buffer [32]byte
	_, err := io.ReadFull(rand.Reader, buffer[:])
	if err != nil {
		log.Fatal(err)
	}
	return buffer
}

func main() {
	nonce := Nonce()
	key := genKey()
	encrypted := secretbox.Seal(nonce[:], []byte("My Secret Message"), &nonce, &key)
	result := hex.EncodeToString(encrypted)
	fmt.Println(result)
}

func generateUserKeys() (*[32]byte, *[32]byte, error) {
	senderPublicKey, senderPrivateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	return senderPublicKey, senderPrivateKey, err
}

func Nonce() [24]byte {
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}
	return nonce
}

func GenerateLazenkey() [32]byte {
	return genKey()
}

func EncryptWithUserKey(publicKey *[32]byte, privateKey *[32]byte, plaintext []byte) *[]byte {
	nonce := Nonce()
	result := box.Seal(nonce[:], plaintext, &nonce, publicKey, privateKey)
	return &result
}

func GenerateUserKeys() (*[32]byte, *[32]byte) {
	senderPublicKey, senderPrivateKey, err := generateUserKeys()
	if err != nil {
		panic(err)
	}
	return senderPublicKey, senderPrivateKey
}

func ToHexString(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

func Chunk(longString string) []string {
	chunkSize := 80 // in bytes
	//longString := "Juni månad blev den varmaste sedan mätningarna började för 135 år sedan, meddelar vetenskapsinstitutet National Oceanic and Atmospheric Administration (NOAA) i USA i sin månadsrapport."
	slices := []string{}
	lastIndex := 0
	lastI := 0
	for i, _ := range longString {
		if i-lastIndex > chunkSize {
			slices = append(slices, longString[lastIndex:lastI])
			lastIndex = lastI
		}
		lastI = i
	}
	// handle the leftovers at the end
	if len(longString)-lastIndex > chunkSize {
		slices = append(slices, longString[lastIndex:lastIndex+chunkSize], longString[lastIndex+chunkSize:])
	} else {
		slices = append(slices, longString[lastIndex:])
	}

	return append(slices, "")
}
