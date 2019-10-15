package main

import (
	"context"
	"fmt"
	"github.com/teamsnap/vault/pkg/vault"
)

var env = map[string]map[string]string{}

var envArr = []string{
	"test/data/test",
}

func main() {
	ctx := context.Background()

	vault.GetSecrets(ctx, &env, envArr)

	fmt.Println("Environment Values:", env)

	fmt.Println("hello = " + env["test/data/test"]["hello"])
}
