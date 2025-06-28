package utils

import (
	"errors"
	"fmt"
	"time"
	"todo-list-api/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	Secret string
}

// Custom Claims structure for better JWT handling
type Claims struct {
	UserID    uint64 `json:"user_id"`
	UserEmail string `json:"user_email"`
	TokenType string `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// TokenPair represents both access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"` // seconds until access token expires
}

// CreateToken creates a token pair (access + refresh) for the user
func (j *JWT) CreateToken(user *models.User) (*models.Token, error) {
	tokenPair, err := j.CreateTokenPair(user)
	if err != nil {
		return nil, err
	}

	// Convert to the existing Token model for backward compatibility
	token := &models.Token{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}

	return token, nil
}

// CreateTokenPair creates both access and refresh tokens
func (j *JWT) CreateTokenPair(user *models.User) (*TokenPair, error) {
	now := time.Now()
	accessTokenExpiry := now.Add(15 * time.Minute)
	refreshTokenExpiry := now.Add(7 * 24 * time.Hour) // 7 days

	// Create Access Token
	accessClaims := &Claims{
		UserID:    user.ID,
		UserEmail: user.Email,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "todo-list-api",
			Subject:   fmt.Sprintf("user:%d", user.ID),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(j.Secret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Create Refresh Token
	refreshClaims := &Claims{
		UserID:    user.ID,
		UserEmail: user.Email,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshTokenExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "todo-list-api",
			Subject:   fmt.Sprintf("user:%d", user.ID),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(j.Secret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessTokenExpiry.Sub(now).Seconds()),
	}, nil
}

// ValidateToken validates a JWT token and returns the claims
func (j *JWT) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Check if token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Additional validation
	if err := j.validateClaims(claims); err != nil {
		return nil, err
	}

	return claims, nil
}

// ValidateAccessToken specifically validates access tokens
func (j *JWT) ValidateAccessToken(tokenString string) (*Claims, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "access" {
		return nil, errors.New("token is not an access token")
	}

	return claims, nil
}

// ValidateRefreshToken specifically validates refresh tokens
func (j *JWT) ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "refresh" {
		return nil, errors.New("token is not a refresh token")
	}

	return claims, nil
}

// RefreshAccessToken creates a new access token using a valid refresh token
func (j *JWT) RefreshAccessToken(refreshTokenString string) (*TokenPair, error) {
	claims, err := j.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Create a new access token
	now := time.Now()
	accessTokenExpiry := now.Add(15 * time.Minute)

	newAccessClaims := &Claims{
		UserID:    claims.UserID,
		UserEmail: claims.UserEmail,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "todo-list-api",
			Subject:   fmt.Sprintf("user:%d", claims.UserID),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newAccessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(j.Secret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign new access token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString, // Keep the same refresh token
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessTokenExpiry.Sub(now).Seconds()),
	}, nil
}

// validateClaims performs additional validation on claims
func (j *JWT) validateClaims(claims *Claims) error {
	// Check if token type is valid
	if claims.TokenType != "access" && claims.TokenType != "refresh" {
		return errors.New("invalid token type")
	}

	// Check if user ID is valid
	if claims.UserID == 0 {
		return errors.New("invalid user ID in token")
	}

	// Check if email is present
	if claims.UserEmail == "" {
		return errors.New("missing user email in token")
	}

	// Check issuer
	if claims.Issuer != "todo-list-api" {
		return errors.New("invalid token issuer")
	}

	return nil
}

// GetUserIDFromToken extracts user ID from a valid token
func (j *JWT) GetUserIDFromToken(tokenString string) (uint64, error) {
	claims, err := j.ValidateAccessToken(tokenString)
	if err != nil {
		return 0, err
	}

	return claims.UserID, nil
}

// GetUserEmailFromToken extracts user email from a valid token
func (j *JWT) GetUserEmailFromToken(tokenString string) (string, error) {
	claims, err := j.ValidateAccessToken(tokenString)
	if err != nil {
		return "", err
	}

	return claims.UserEmail, nil
}
