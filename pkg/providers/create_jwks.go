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
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	awsProvider "github.com/chuhaoyuu/aws-oidc-sts/pkg/providers/aws"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

// CreateIdentityProvider creates an identity provider by generating a JSON Web Key Set (JWKS), creating a signed JWT, and setting up an S3 bucket.
// Parameters:
// - filePath: Path to the private key file used for creating the JWKS.
// - bucketName: Name of the S3 bucket to be created.
// - region: AWS region where the S3 bucket will be created.
// Returns:
// - error: An error if any step fails during the creation process.
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

	identity, err := client.AwsClientIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return fmt.Errorf("failed to get AWS client identity: %w", err)
	}
	slog.Info("AWS Client Identity", "Account", aws.ToString(identity.Account), "Arn", aws.ToString(identity.Arn), "Region", region, "UserId", aws.ToString(identity.UserId))

	slog.Info("Creating S3 bucket", "BucketName", bucketName)
	buckerInfo, err := client.CreateS3Bucket(&s3.CreateBucketInput{Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create S3 bucket: %w", err)
	}
	slog.Info("S3 Bucket created successfully", "Bucket Location", aws.ToString(buckerInfo.Location))

	return nil
}

// keyIDFromPublicKey generates a unique key identifier (key ID) from the given public key.
// Parameters:
// - publicKey: The public key to generate the key ID from.
// Returns:
// - string: A unique key ID derived from the public key.
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

// CreateJSONWebKeySet creates a JSON Web Key Set (JWKS) from the provided private key file.
// Parameters:
// - filePath: Path to the private key file used for creating the JWKS.
// Returns:
// - jwk.Key: The private key imported as a JWK.
// - error: An error if any step fails during the creation process.
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
