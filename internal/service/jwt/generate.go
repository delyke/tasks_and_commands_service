package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// Generate creates a new JWT token.
func (g *Generator) Generate(userID user.ID, email string, ttl time.Duration, issuer string) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    issuer,
		},
		UserID: userID,
		Email:  email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(g.secretKey)
}
