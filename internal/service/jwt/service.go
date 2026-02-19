package jwt

import (
	"github.com/golang-jwt/jwt/v5"

	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// Generator implements TokenGenerator using JWT.
type Generator struct {
	secretKey []byte
}

// NewJWTGenerator creates a new JWTGenerator.
func NewJWTGenerator(secretKey string) *Generator {
	return &Generator{
		secretKey: []byte(secretKey),
	}
}

// Claims represents JWT claims.
type Claims struct {
	jwt.RegisteredClaims
	UserID user.ID `json:"user_id"`
	Email  string  `json:"email"`
}
