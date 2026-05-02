package jwths256

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	ErrMalformed = errors.New("jwt malformed")
	ErrInvalid   = errors.New("jwt invalid")
	ErrExpired   = errors.New("jwt expired")
)

// Claims is the minimal JWT payload used for access/refresh tokens.
type Claims struct {
	Typ string `json:"typ"`
	Sub string `json:"sub"`
	Jti string `json:"jti"`
	Exp int64  `json:"exp"`
	Iat int64  `json:"iat"`
	Nbf int64  `json:"nbf"`
}

func b64URL(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func b64URLDecode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

// Sign builds a compact HS256 JWT.
func Sign(secret []byte, c Claims) (string, error) {
	hdr := map[string]string{"alg": "HS256", "typ": "JWT"}
	hb, err := json.Marshal(hdr)
	if err != nil {
		return "", err
	}
	pb, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	head := b64URL(hb)
	payload := b64URL(pb)
	unsigned := head + "." + payload
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(unsigned))
	sig := b64URL(mac.Sum(nil))
	return unsigned + "." + sig, nil
}

// Parse verifies HS256 signature and unmarshals claims.
// If validateExp is false, Exp is not checked (logout path for stale refresh).
func Parse(secret []byte, token string, validateExp bool) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrMalformed
	}
	unsigned := parts[0] + "." + parts[1]
	sig, err := b64URLDecode(parts[2])
	if err != nil {
		return nil, ErrMalformed
	}
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(unsigned))
	if !hmac.Equal(sig, mac.Sum(nil)) {
		return nil, ErrInvalid
	}
	pb, err := b64URLDecode(parts[1])
	if err != nil {
		return nil, ErrMalformed
	}
	var c Claims
	if err := json.Unmarshal(pb, &c); err != nil {
		return nil, ErrMalformed
	}
	now := time.Now().Unix()
	if validateExp {
		if c.Exp > 0 && now >= c.Exp {
			return nil, ErrExpired
		}
		if c.Nbf > 0 && now < c.Nbf {
			return nil, ErrInvalid
		}
	}
	return &c, nil
}
