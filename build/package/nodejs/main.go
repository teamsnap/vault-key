package main

import (
	"C"
	"fmt"

	"github.com/teamsnap/vault/pkg/vault/key"
)

var env = map[string]map[string]string{}

//export GetSecrets
func GetSecrets(secretNames *C.char) (*C.char, *C.char) {
	secretNamesStr := C.GoString(secretNames)

	secrets, err := key.Loot(secretNamesStr)
	if err != nil {
		return C.CString(""), C.CString(fmt.Sprintf("Error: ", err))
	}

	return C.CString(string(secrets)), C.CString("")
}

func main() {}
