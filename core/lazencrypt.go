package core

import (
	"io"
	"crypto/rand"
	"log"
	"encoding/hex"
	"golang.org/x/crypto/nacl/box"
	"io/ioutil"
	"path"
	"os/user"
	"fmt"
	"golang.org/x/crypto/nacl/secretbox"
	"github.com/TwistedSystemsLtd/lazenby-go/lazendata"
)

func genKey() [32]byte {
	var buffer [32]byte
	_, err := io.ReadFull(rand.Reader, buffer[:])
	if err != nil {
		log.Fatal(err)
	}
	return buffer
}

func generateUserKeys() (*[32]byte, *[32]byte, error) {
	senderPublicKey, senderPrivateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	return senderPublicKey, senderPrivateKey, err
}

func ReadUserKeys(userKeyDir string) (*[32]byte, *[32]byte) {
	publicKey, publicErr := ioutil.ReadFile(path.Join(userKeyDir, "publicKey"))
	privateKey, privateErr := ioutil.ReadFile(path.Join(userKeyDir, "privateKey"))

	if publicErr != nil || privateErr != nil {
		log.Panic("Could not read user keys", userKeyDir)
	}

	var publicBytes [32]byte
	var privateBytes [32]byte

	_, pubDecodeErr := hex.Decode(publicBytes[:], publicKey)
	_, prvDecodeErr := hex.Decode(privateBytes[:], privateKey)

	if pubDecodeErr != nil || prvDecodeErr != nil {
		log.Panic("Could not decode user key data")
	}

	return &publicBytes, &privateBytes
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

func EncryptWithUserKey(publicKey *[32]byte, privateKey *[32]byte, plaintext []byte) []byte {
	fmt.Println("Plaintext", ToHexString(plaintext))
	nonce := Nonce()
	return box.Seal(nonce[:], plaintext, &nonce, publicKey, privateKey)
}

func DecryptWithUserKey(publicKey *[32]byte, privateKey *[32]byte, encrypted []byte) []byte {
	var decryptNonce [24]byte
	copy(decryptNonce[:], encrypted[:24])
	decrypted, ok := box.Open(nil, encrypted[24:], &decryptNonce, publicKey, privateKey)
	if !ok {
		panic("decryption error")
	}

	fmt.Println("Decrypted", ToHexString(decrypted))

	return decrypted
}

func EncryptWithLazenkey(lazenkey *[32]byte, plaintext []byte) []byte {
	nonce := Nonce()
	return secretbox.Seal(nonce[:], plaintext, &nonce, lazenkey)
}

func DecryptWithLazenkey(lazenkey *[32]byte, encrypted []byte) []byte {
	var decryptNonce [24]byte
	copy(decryptNonce[:], encrypted[:24])
	decrypted, ok := secretbox.Open(nil, encrypted[24:], &decryptNonce, lazenkey)
	if !ok {
		panic("decryption error")
	}
	return decrypted
}

func DecryptLazenkey(publicKey *[32]byte, privateKey *[32]byte, lazendata *lazendata.Lazenfile) *[32]byte {
	usersEncryptedLazenkey := lazendata.Lazenkeys[ToHexString(publicKey[:])]
	var lazenkey [32]byte

	result := DecryptWithUserKey(publicKey, privateKey, usersEncryptedLazenkey.Lazenkey)

	copy(lazenkey[:], result[:32])
	return &lazenkey
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

func FromHexString(hexString string) []byte {
	result, err := hex.DecodeString(hexString)
	if err != nil {
		log.Panic("Error decoding hex string")
	}
	return result
}

func Chunk(longString string) []string {
	chunkSize := 80 // in bytes
	var slices []string
	lastIndex := 0
	lastI := 0
	for i := range longString {
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

func Lazenhome() string {
	currentUser := CurrentUser()

	home := currentUser.HomeDir
	return path.Join(home, ".lzb")
}

func CurrentUser() *user.User {
	currentUser, err := user.Current()
	if err != nil {
		log.Panic(err)
	}
	return currentUser
}
