package auth

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

	"payrool/internal/user"
)

const DefaultSecret = "payrool-dev-secret"

const (
	RoleAdmin             = "admin"
	RoleManager           = "manager"
	RoleAttendanceOfficer = "attendance_officer"
	RoleStaff             = "staff"
)

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Exp      int64  `json:"exp"`
}

func secret() string {
	if s := strings.TrimSpace(os.Getenv("PAYROOL_JWT_SECRET")); s != "" {
		return s
	}
	return DefaultSecret
}

func GenerateToken(u user.User) (string, error) {
	claims := Claims{
		UserID:   u.ID,
		Username: u.Username,
		Role:     u.Role,
		Exp:      time.Now().Add(8 * time.Hour).Unix(),
	}

	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	signature := hmac.New(sha256.New, []byte(secret()))
	_, err = signature.Write([]byte(encodedPayload))
	if err != nil {
		return "", err
	}

	mac := base64.RawURLEncoding.EncodeToString(signature.Sum(nil))
	return encodedPayload + "." + mac, nil
}

func ParseToken(token string) (Claims, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return Claims{}, errors.New("missing token")
	}

	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return Claims{}, errors.New("invalid token format")
	}

	encodedPayload := parts[0]
	providedMAC := parts[1]

	signature := hmac.New(sha256.New, []byte(secret()))
	_, err := signature.Write([]byte(encodedPayload))
	if err != nil {
		return Claims{}, err
	}

	expectedMAC := base64.RawURLEncoding.EncodeToString(signature.Sum(nil))
	if !hmac.Equal([]byte(expectedMAC), []byte(providedMAC)) {
		return Claims{}, errors.New("invalid token signature")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(encodedPayload)
	if err != nil {
		return Claims{}, err
	}

	var claims Claims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return Claims{}, err
	}

	if claims.Exp == 0 || time.Now().Unix() > claims.Exp {
		return Claims{}, errors.New("token expired")
	}

	return claims, nil
}

func RequireRole(requiredRoles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authorization := strings.TrimSpace(r.Header.Get("Authorization"))
			if authorization == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authorization, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			claims, err := ParseToken(parts[1])
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if len(requiredRoles) > 0 {
				allowed := false
				for _, required := range requiredRoles {
					if strings.EqualFold(claims.Role, required) {
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
