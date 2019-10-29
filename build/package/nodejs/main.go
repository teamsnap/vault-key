package main

import (
	"C"

	"github.com/google/martian/log"
	"github.com/teamsnap/vault/pkg/vault"
)

//export GetSecrets
func GetSecrets(secretNames *C.char) *C.char {
	secretNamesStr := C.GoString(secretNames)

	secrets, err := vault.Loot(secretNamesStr)
	if err != nil {
		log.Fatal(err)
	}

	return C.CString(string(secrets))
}

func main() {}
