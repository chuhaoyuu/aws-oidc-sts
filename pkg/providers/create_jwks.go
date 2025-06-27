package providers

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsProvider "github.com/chuhaoyuu/aws-oidc-sts/pkg/providers/aws"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

func CreateIdentityProvider(filePath, bucketName, region string) error {
	// Create the JWKS file
	jwkKey, err := CreateJSONWebKeySet(filePath)
	if err != nil {
		return fmt.Errorf("failed to create JSON Web Key Set: %w", err)
	}

	signedJWT, err := CreateJWT(jwkKey)
	if err != nil {
		return fmt.Errorf("failed to create JWT: %w", err)
	}

	slog.Info("JWT created successfully", "JWT", string(signedJWT))

	client, err := awsProvider.NewAwsFromConfig(region)
	if err != nil {
		return fmt.Errorf("failed to create AWS client: %w", err)
	}

	identity, err := client.AwsClientIdentity()
	if err != nil {
		return fmt.Errorf("failed to get AWS client identity: %w", err)
	}

	slog.Info("AWS Client Identity",
		"Account", aws.ToString(identity.Account),
		"Arn", aws.ToString(identity.Arn),
		"Region", region,
		"UserId", aws.ToString(identity.UserId),
	)
	err = client.CreateS3Bucket(bucketName, region)
	if err != nil {
		return fmt.Errorf("failed to create S3 bucket: %w", err)
	}

	return nil
}

// keyIDFromPublicKey generates a unique key identifier (key ID) from the given public key.
// The publicKey parameter can be of any type that represents a public key.
// This function is typically used to create a key ID for use in JSON Web Key Sets (JWKS).
// The returned key ID is a string that uniquely identifies the provided public key.
func keyIDFromPublicKey(publicKey any) string {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		panic(fmt.Errorf("failed to marshal public key: %w", err))
	}

	hash := sha256.New()
	hash.Write(publicKeyBytes)
	hashedBytes := hash.Sum(nil)

	keyID := hex.EncodeToString(hashedBytes)

	return keyID
}

// CreateJSONWebKeySet generates a JSON Web Key Set (JWKS) from a given private key file.
// It parses the private and public keys from the specified file, creates a JWK for the private key,
// and adds the corresponding public key to the JWK Set. The resulting JWKS is then written to a file.
//
// Parameters:
//   - filePath: The path to the private key file.
//
// Returns:
//   - jwk.Key: The generated JWK for the private key.
//   - error: An error if any step in the process fails.
//
// The function performs the following steps:
//  1. Parses the private and public keys from the provided file.
//  2. Creates a new JWK Set.
//  3. Extracts the key ID (kid) from the public key.
//  4. Imports the private key into a JWK and sets its key ID, usage, and algorithm.
//  5. Extracts the public key from the private key and adds it to the JWK Set.
//  6. Marshals the JWK Set into JSON format.
//  7. Writes the JSON-formatted JWK Set to a file in the specified directory.
//
// Errors are returned if any of the following occur:
//   - Parsing the private or public key fails.
//   - Importing the private key into a JWK fails.
//   - Setting the key ID, usage, or algorithm for the JWK fails.
//   - Creating the public key from the private key fails.
//   - Marshaling the JWK Set into JSON format fails.
//   - Writing the JWK Set to a file fails.
func CreateJSONWebKeySet(filePath string) (jwk.Key, error) {

	privateKey, err := ParsePrivateKeyFromFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	publicKey, err := ParsePublicKeyFromFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Create a new JWK Set
	jwkSet := jwk.NewSet()

	// Extract the key ID (kid) from the public key
	keyID := keyIDFromPublicKey(publicKey)

	// Import the RSA private key into a JWK
	jwkPrivateKey, err := jwk.Import(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to import private key into JWK: %w", err)
	}

	// Set the key ID (kid) for the JWK
	if err := jwkPrivateKey.Set(jwk.KeyIDKey, keyID); err != nil {
		return nil, fmt.Errorf("failed to set key ID: %w", err)
	}

	// Set the key type (use) for the JWK
	if err := jwkPrivateKey.Set(jwk.KeyUsageKey, JWKSUsage); err != nil {
		return nil, fmt.Errorf("failed to set key usage: %w", err)
	}

	// Set the algorithm (alg) for the JWK
	if err := jwkPrivateKey.Set(jwk.AlgorithmKey, jwa.RS256()); err != nil {
		return nil, fmt.Errorf("failed to set algorithm: %w", err)
	}

	// Extract the public key from the private key
	jwkPublicKey, err := jwk.PublicKeyOf(jwkPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create public key from private key: %w", err)
	}

	// Add the public key to the JWK Set
	jwkSet.AddKey(jwkPublicKey)

	// Marshal the JWK Set into JSON format
	jwkSetJSON, err := json.MarshalIndent(jwkSet, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JWK Set: %w", err)
	}

	// Write the JWK Set to a file
	jwkFilePath := filepath.Join(filePath, TLSDirName, JWKSFileName)
	if err := os.WriteFile(jwkFilePath, jwkSetJSON, 0644); err != nil {
		return nil, fmt.Errorf("failed to write JWK Set to file: %w", err)
	}

	return jwkPrivateKey, nil
}
