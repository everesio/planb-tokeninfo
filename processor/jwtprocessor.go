package processor

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtProcessor interface {
	Process(t *jwt.Token, timeBase time.Time) (*TokenInfo, error)
}

// TokenInfo type is used to serialize a JWT validation result in a standard Token Info JSON format
type TokenInfo struct {
	AccessToken   string            `json:"access_token"`
	RefreshToken  string            `json:"refresh_token,omitempty"`
	UID           string            `json:"uid"`
	GrantType     string            `json:"grant_type"`
	Scope         []string          `json:"scope"`
	Realm         string            `json:"realm"`
	ClientId      string            `json:"client_id"`
	TokenType     string            `json:"token_type"`
	ExpiresIn     int               `json:"expires_in"`
	PrivateClaims map[string]string `json:"-"`
}
