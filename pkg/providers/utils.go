package providers

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
)

// ParsePublicKeyFromFile reads a PEM-encoded RSA public key from a specified file,
// decodes the PEM block, and parses the public key.
//
// Parameters:
//   - filePath: The path to the directory containing the public key file.
//
// Returns:
//   - any: The parsed public key object.
//   - error: An error if the file cannot be read, the PEM block cannot be decoded,
//     or the public key cannot be parsed.
//
// Errors:
//   - Returns an error if the file cannot be read.
//   - Returns an error if the PEM block is invalid or cannot be decoded.
//   - Returns an error if the public key cannot be parsed.
func ParsePublicKeyFromFile(filePath string) (any, error) {
	// Read the public key from the specified file
	publicKeyPem, err := os.ReadFile(filepath.Join(filePath, TLSDirName, RSAPublicKeyFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	// Decode the PEM-encoded public key
	block, _ := pem.Decode(publicKeyPem)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	// Parse the public key
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return publicKey, nil
}

// ParsePrivateKeyFromFile reads an RSA private key from a specified file path,
// decodes the PEM-encoded private key, and parses it into an *rsa.PrivateKey.
//
// Parameters:
//   - filePath: The path to the directory containing the private key file.
//
// Returns:
//   - *rsa.PrivateKey: The parsed RSA private key.
//   - error: An error if the file cannot be read, the PEM block cannot be decoded,
//     or the private key cannot be parsed.
//
// The function expects the private key file to be named as specified by the
// RSAPrivateKeyFile constant and located in the provided directory path.
func ParsePrivateKeyFromFile(filePath string) (*rsa.PrivateKey, error) {
	// Read the private key from the specified file
	privateKeyPem, err := os.ReadFile(filepath.Join(filePath, TLSDirName, RSAPrivateKeyFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	// Decode the PEM-encoded private key
	block, _ := pem.Decode(privateKeyPem)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	// Parse the private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privateKey, nil
}
