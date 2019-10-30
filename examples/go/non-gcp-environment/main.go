package main

import (
	"context"
	"fmt"
	"log"

	"github.com/teamsnap/vault-key/pkg/vault"
)

var env = map[string]map[string]string{}

var envArr = []string{
	"test/data/test",
}

func main() {
	ctx := context.Background()

	err := vault.GetSecrets(ctx, &env, envArr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Environment Values:", env)

	fmt.Println("hello = " + env["test/data/test"]["hello"])
}
