// package name: vault
package main

import "C"
import (
	"context"
	"github.com/stevenaldinger/vault/pkg/vault"
	"fmt"
	"encoding/json"
)

var env = map[string]map[string]string{}

//export GetSecrets
func GetSecrets(secretNames *C.char) *C.char {
		ctx := context.Background()

		secretNamesStr := C.GoString(secretNames)

		var envArr []string
		if err := json.Unmarshal([]byte(secretNamesStr), &envArr); err != nil {
			panic(err)
		}

		vault.GetSecrets(ctx, &env, envArr)

		secretData, err := json.Marshal(env)
		if err != nil {
				fmt.Println(err.Error())
				return C.CString("")
		}

		jsonStr := string(secretData)

		return C.CString(jsonStr)
}

func main() {
}
