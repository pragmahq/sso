package web

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// The secret is a generated 32-bit string used to provide a signature to our signed tokens.
var SECRET = []byte(os.Getenv("COOKIE_SECRET"))

// Claims represents the custom claims structure for JWT tokens.
//
// It includes the email address and expiration time (exp) as specified in the
// JSON Web Token (JWT) standard, along with additional registered claims.
type Claims struct {
	Email                string `json:"email"` // Email address associated with the token.
	Exp                  int    `json:"exp"`   // Expiration time of the token (Unix timestamp).
	jwt.RegisteredClaims        // Embedded struct for standard JWT claims.
}

// CreateToken generates a new JWT token with the provided email and a default expiration time.
//
// Parameters:
//   - email: The email address to be included in the token claims.
//
// Returns:
//   - string: The generated JWT token as a string.
//   - error: An error, if any, encountered during token generation.
func CreateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 720).Unix(),
	})

	tokenStr, err := token.SignedString(SECRET)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// DecodeToken decodes a JWT token and returns the custom claims if the token is valid.
//
// Parameters:
//   - tokenString: The JWT token to be decoded.
//
// Returns:
//   - *Claims: The custom claims decoded from the token.
//   - error: An error, if any, encountered during token decoding.
func DecodeToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return SECRET, nil
	}, jwt.WithLeeway(5*time.Second))

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
