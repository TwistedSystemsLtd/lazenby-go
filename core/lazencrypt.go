package core

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"github.com/TwistedSystemsLtd/lazenby-go/lazendata"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/nacl/secretbox"
	"io"
	"io/ioutil"
	"log"
	"os/user"
	"path"
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

	_, pubDecodeErr := base64.RawURLEncoding.Decode(publicBytes[:], publicKey)
	_, prvDecodeErr := base64.RawURLEncoding.Decode(privateBytes[:], privateKey)

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
	usersEncryptedLazenkey := lazendata.Lazenkeys[EncodeString(publicKey[:])]
	var lazenkey [32]byte

	result := DecryptWithUserKey(publicKey, privateKey, DecodeString(usersEncryptedLazenkey))

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

func EncodeString(bytes []byte) string {
	return base64.URLEncoding.EncodeToString(bytes)
}

func DecodeString(base64String string) []byte {
	result, err := base64.URLEncoding.DecodeString(base64String)
	if err != nil {
		log.Panic("Error decoding base64 string")
	}
	return result
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
