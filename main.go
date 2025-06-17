package main

import (
	"fmt"
	"os"

	"github.com/chuhaoyuu/aws-oidc-sts/pkg/providers"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current working directory: %v\n", err)
		os.Exit(1)
	}
	if err := providers.CreateRSAKeyPair(pwd); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating RSA key pair: %v\n", err)
		os.Exit(1)
	}
}
