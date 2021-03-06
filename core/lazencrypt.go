// Lazenby - Your secrets as a service
// Copyright © 2018 Twisted Systems Ltd
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
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
	publicKeyByteString, publicErr := ioutil.ReadFile(path.Join(userKeyDir, "publicKey"))
	privateKeyByteString, privateErr := ioutil.ReadFile(path.Join(userKeyDir, "privateKey"))

	if publicErr != nil || privateErr != nil {
		log.Panic("Could not read user keys", userKeyDir)
	}

	return DecodeKeyStrings(string(publicKeyByteString), string(privateKeyByteString))
}

func DecodeKeyStrings(publicKey string, privateKey string) (*[32]byte, *[32]byte) {
	publicBytes := DecodeString(publicKey)
	privateBytes := DecodeString(privateKey)

	return ToKeyByteArray(publicBytes), ToKeyByteArray(privateBytes)
}

func ToKeyByteArray(keySlice []byte) *[32]byte {
	var keyBytes [32]byte
	copy(keyBytes[:], keySlice)
	return &keyBytes
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

func EncryptWithUserKey(publicEncryptionKey *[32]byte, privateSigningKey *[32]byte, publicSigningKey *[32] byte, plaintext []byte) []byte {
	nonce := Nonce()
	preamble := append(nonce[:], publicSigningKey[:]...)

	println(len(preamble))

	return box.Seal(preamble, plaintext, &nonce, publicEncryptionKey, privateSigningKey)
}

func DecryptWithUserKey(privateDecryptionKey *[32]byte, encrypted []byte) []byte {
	log.Println(len(encrypted))
	log.Println(EncodeString(encrypted))
	var decryptNonce [24]byte
	var publicSigningKey [32]byte

	log.Println(len(encrypted[:24]))
	copy(decryptNonce[:], encrypted[:24])

	fmt.Println(len(encrypted[24:(32+24)]))
	copy(publicSigningKey[:], encrypted[24:(32+24)])

	decrypted, ok := box.Open(nil, encrypted[(32+24):], &decryptNonce, &publicSigningKey, privateDecryptionKey)
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

	log.Println("Users lazenkey", usersEncryptedLazenkey)


	var lazenkey [32]byte

	result := DecryptWithUserKey(privateKey, DecodeString(usersEncryptedLazenkey))

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
