package main

import (
	"fmt"
	"log/slog"
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

	jwkKey, err := providers.CreateJSONWebKeySet(pwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating JSON Web Key Set: %v\n", err)
		os.Exit(1)
	}
	signedJWT, err := providers.CreateJWT(jwkKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating JWT: %v\n", err)
		os.Exit(1)
	}
	slog.Info("Successfully created RSA key pair, JWKS, and JWT.")
	fmt.Printf("Signed JWT: %s\n", signedJWT)

}
