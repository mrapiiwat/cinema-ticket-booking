package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Exp    int64  `json:"exp"`
}

func GenerateToken(userID, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		Exp:    time.Now().Add(24 * time.Hour).Unix(),
	}

	header, err := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	if err != nil {
		return "", err
	}
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(header)
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	unsigned := encodedHeader + "." + encodedPayload
	signature := sign(unsigned)

	return unsigned + "." + signature, nil
}

func ParseToken(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	unsigned := parts[0] + "." + parts[1]
	expected := sign(unsigned)
	if !hmac.Equal([]byte(expected), []byte(parts[2])) {
		return nil, errors.New("invalid token signature")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}
	if claims.UserID == "" || claims.Role == "" {
		return nil, errors.New("invalid token claims")
	}
	if time.Now().Unix() > claims.Exp {
		return nil, errors.New("token expired")
	}

	return &claims, nil
}

func sign(unsigned string) string {
	mac := hmac.New(sha256.New, []byte(secret()))
	mac.Write([]byte(unsigned))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func secret() string {
	value := os.Getenv("JWT_SECRET")
	if value == "" {
		return "cinema-ticket-booking-dev-secret"
	}
	return value
}
