package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yeying-community/webdav/internal/domain/auth"
)

// JWTManager JWT 管理器
type JWTManager struct {
	secret     []byte
	expiration time.Duration
	issuer     string
}

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

// Claims JWT 声明
type Claims struct {
	Address   string `json:"address"`
	TokenType string `json:"token_type,omitempty"`
	jwt.RegisteredClaims
}

// NewJWTManager 创建 JWT 管理器
func NewJWTManager(secret string, expiration time.Duration) *JWTManager {
	return &JWTManager{
		secret:     []byte(secret),
		expiration: expiration,
		issuer:     "webdav-server",
	}
}

// Generate 生成 JWT
func (m *JWTManager) Generate(address string) (*auth.Token, error) {
	return m.generate(address, TokenTypeAccess, m.expiration)
}

func (m *JWTManager) GenerateRefresh(address string, expiration time.Duration) (*auth.Token, error) {
	return m.generate(address, TokenTypeRefresh, expiration)
}

func (m *JWTManager) generate(address, tokenType string, expiration time.Duration) (*auth.Token, error) {
	now := time.Now()
	expiresAt := now.Add(expiration)

	claims := Claims{
		Address:   strings.ToLower(address),
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    m.issuer,
			Subject:   address,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return &auth.Token{
		Value:     tokenString,
		Address:   address,
		ExpiresAt: expiresAt,
		IssuedAt:  now,
	}, nil
}

// Verify 验证 JWT
func (m *JWTManager) Verify(tokenString string) (string, error) {
	return m.verify(tokenString, TokenTypeAccess, true)
}

func (m *JWTManager) VerifyRefresh(tokenString string) (string, error) {
	return m.verify(tokenString, TokenTypeRefresh, false)
}

func (m *JWTManager) verify(tokenString, expectedType string, allowEmptyType bool) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", auth.ErrTokenExpired
		}
		return "", fmt.Errorf("%w: %v", auth.ErrInvalidToken, err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if claims.TokenType == "" && allowEmptyType {
			return claims.Address, nil
		}
		if claims.TokenType != expectedType {
			return "", auth.ErrInvalidToken
		}
		return claims.Address, nil
	}

	return "", auth.ErrInvalidToken
}
