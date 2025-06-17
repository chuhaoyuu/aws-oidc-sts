package providers

// Creates a new RSA key pair for use with AWS OIDC STS
import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

// CreateRSAKeyPair generates an RSA key pair (private and public keys) and saves them to the specified file path.
// If the key pair already exists at the specified location, the function skips the generation process.
//
// Parameters:
//   - keyPairFilePath: The base directory where the RSA key pair will be stored. The private key will be saved
//     as a file named "RSAPrivateKeyFile" and the public key as "RSAPublicKeyFile" within a subdirectory.
//
// Behavior:
//   - Creates the necessary directory structure if it does not exist.
//   - Checks if the private and public key files already exist. If both files are present, the function logs
//     a warning and skips the key generation process.
//   - If the key pair does not exist, generates a 4096-bit RSA private key and derives the public key from it.
//   - Encodes the private key in PEM format and writes it to the private key file with restricted permissions (0600).
//   - Encodes the public key in PEM format and writes it to the public key file with read permissions (0644).
//
// Returns:
//   - An error if any step in the process fails, such as directory creation, key generation, or file writing.
//   - nil if the key pair is successfully created or already exists.
//
// Logging:
//   - Logs informational messages during the process, including warnings if the key pair already exists.
//   - Logs debug messages for successful file writes.
func CreateRSAKeyPair(keyPairFilePath string) error {

	RSAKeyDir := filepath.Join(keyPairFilePath, TLSDirName)
	if err := os.MkdirAll(RSAKeyDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory for RSA keys: %w", err)
	}
	privateKeyFile := filepath.Join(RSAKeyDir, RSAPrivateKeyFile)
	publicKeyFile := filepath.Join(RSAKeyDir, RSAPublicKeyFile)
	skipGeneration := false

	if _, err := os.Stat(privateKeyFile); err == nil {
		slog.Warn("Private key file already exists, skipping creation.", slog.String("file", privateKeyFile))
		skipGeneration = true
	}
	if _, err := os.Stat(publicKeyFile); err == nil {
		slog.Warn("Public key file already exists, skipping creation.", slog.String("file", publicKeyFile))
		skipGeneration = true
	}

	if skipGeneration {
		slog.Info("RSA key pair already exists, skipping creation.")
		return nil
	}

	// Generate RSA private key
	slog.Info("Generating RSA key pair...")

	bitSize := 4096
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return fmt.Errorf("failed to generate RSA private key: %w", err)
	}

	// Encode private key to PEM format
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// Extract public key from private key
	publicKey := &privateKey.PublicKey

	// Encode public key to PEM format
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %w", err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	// Write private key to file
	slog.Info("Writing private key to", slog.String("file", privateKeyFile))
	err = os.WriteFile(privateKeyFile, privateKeyPEM, 0600)
	if err != nil {
		return fmt.Errorf("failed to write private key to file: %w", err)
	}
	slog.Debug("Private key written successfully", slog.String("file", privateKeyFile))

	// Write public key to file
	slog.Info("Writing public key to", slog.String("file", publicKeyFile))
	err = os.WriteFile(publicKeyFile, publicKeyPEM, 0644)
	if err != nil {
		return fmt.Errorf("failed to write public key to file: %w", err)
	}
	slog.Info("RSA key pair generated successfully.")

	return nil
}
