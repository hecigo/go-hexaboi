package jwt

import (
	"fmt"
	"strings"

	"github.com/goccy/go-json"

	"github.com/golang-jwt/jwt/v4"
)

func ValidateJWTToken[T jwt.Claims](base64Token string) (claims *T, err error) {
	if strings.Contains(base64Token, " ") {
		base64Token = strings.TrimSpace(strings.Split(base64Token, " ")[1])
	}

	jwks := JWKS()

	// Parse token with claims
	token, err := jwt.ParseWithClaims(base64Token, &jwt.MapClaims{}, jwks.Keyfunc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	bytes, err := json.Marshal(token.Claims)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal claims: %w", err)
	}

	if err := json.Unmarshal(bytes, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	return claims, nil
}
