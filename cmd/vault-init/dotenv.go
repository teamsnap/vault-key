package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

// DotEnvVariables represents the variables passed in to create the .env file
type DotEnvVariables struct {
	Secrets   map[string]string
}

const dotEnv = `
{{ range $key, $value := .Secrets}}
export {{ $key }}={{ $value }}
{{end}}
`

// GenerateDotEnv parses the .env template with passed in variables
// and returns a string
func GenerateDotEnv(dotenvVars DotEnvVariables) string {
	t, err := template.New("dotenv").Parse(inventory)
	if err != nil {
		panic(err)
	}

	var tpl bytes.Buffer
	err = t.Execute(&tpl, dotenvVars)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing dotenv template:", err)
		os.Exit(1)
	}

	result := tpl.String()

	return result
}
