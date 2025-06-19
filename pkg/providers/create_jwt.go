package providers

//
// Creates a new JWT for use with AWS OIDC STS
import (
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

// CreateJWT generates a signed JWT token using the provided private key.
//
// The JWT token includes the following claims:
// - "iss" (Issuer): The entity that issued the JWT, defined by the constant JWTIssuer.
// - "aud" (Audience): The intended audience for the JWT, defined by the constant JWTAudience.
// - "sub" (Subject): The subject of the JWT, defined by the constant JWTSubject.
// - "exp" (Expiration Time): The expiration time of the JWT, set to 24 hours from the current time.
// - "iat" (Issued At): The time at which the JWT was issued, set to the current time.
//
// Parameters:
// - signingKey (jwk.Key): The private key used to sign the JWT.
//
// Returns:
// - ([]byte): The signed JWT token as a byte slice.
// - (error): An error if the token creation or signing process fails.
func CreateJWT(signingKey jwk.Key) ([]byte, error) {

	// Create a new JWT token with the specified claims
	token, err := jwt.NewBuilder().Claim("iss", JWTIssuer).
		Claim("aud", JWTAudience).
		Claim("sub", JWTSubject).
		Claim("exp", time.Now().Add(time.Hour*24).Unix()).
		Claim("iat", time.Now().Unix()).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT token: %w", err)
	}

	// Sign the JWT token using the private key
	signedJWT, err := jwt.Sign(token, jwt.WithKey(jwa.RS256(), signingKey))
	if err != nil {
		return nil, fmt.Errorf("failed to sign JWT token: %w", err)
	}

	return signedJWT, nil
}
