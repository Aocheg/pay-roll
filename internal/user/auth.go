package user

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Exp      int64  `json:"exp"`
}

func secret() string {
	if value := strings.TrimSpace(os.Getenv("PAYROOL_JWT_SECRET")); value != "" {
		return value
	}
	return "payrool-dev-secret"
}

func GenerateToken(u User) (string, error) {
	claims := Claims{UserID: u.ID, Username: u.Username, Role: u.Role, Exp: time.Now().Add(8 * time.Hour).Unix()}
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	encoded := base64.RawURLEncoding.EncodeToString(payload)
	mac := hmac.New(sha256.New, []byte(secret()))
	if _, err := mac.Write([]byte(encoded)); err != nil {
		return "", err
	}
	return encoded + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), nil
}

func ParseToken(token string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return Claims{}, errors.New("invalid token format")
	}
	mac := hmac.New(sha256.New, []byte(secret()))
	if _, err := mac.Write([]byte(parts[0])); err != nil {
		return Claims{}, err
	}
	provided, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil || !hmac.Equal(mac.Sum(nil), provided) {
		return Claims{}, errors.New("invalid token signature")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Claims{}, errors.New("invalid token payload")
	}
	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return Claims{}, errors.New("invalid token payload")
	}
	if claims.Exp == 0 || time.Now().Unix() > claims.Exp {
		return Claims{}, errors.New("token expired")
	}
	return claims, nil
}

func RequireRole(requiredRoles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			parts := strings.SplitN(strings.TrimSpace(r.Header.Get("Authorization")), " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}
			claims, err := ParseToken(strings.TrimSpace(parts[1]))
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			if len(requiredRoles) > 0 {
				allowed := false
				for _, role := range requiredRoles {
					if strings.EqualFold(claims.Role, role) {
						allowed = true
						break
					}
				}
				if !allowed {
					http.Error(w, fmt.Sprintf("requires one of: %s", strings.Join(requiredRoles, ", ")), http.StatusForbidden)
					return
				}
			}
			next(w, r)
		}
	}
}
