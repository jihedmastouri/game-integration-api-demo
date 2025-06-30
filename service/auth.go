package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jihedmastouri/game-integration-api-demo/internal"
	"github.com/jihedmastouri/game-integration-api-demo/models"
	"golang.org/x/crypto/bcrypt"
)

type ClaimType struct {
	jwt.RegisteredClaims
	SessionID string `json:"session,omitempty"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Service) AuthenticatePlayer(ctx context.Context, req AuthRequest) (token string, err error) {
	player, err := s.Repository.GetPlayerByUsername(ctx, req.Username)
	if err != nil {
		return "", errors.New("USER NOT FOUND")
	}

	if !s.validatePassword(req.Password, player.Password) {
		return "", errors.New("PASSWORD MISMATCH")
	}

	playerSession, err := s.Repository.CreatePlayerSession(ctx, player.ID)
	if err != nil {
		return "", err
	}

	return s.generateJWT(playerSession)
}

func (s *Service) AuthorizePlayer(ctx context.Context, token string) (*models.Player, error) {
	claims, err := s.validateJWT(token)
	if err != nil {
		return nil, err
	}

	numDate, err := claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}

	if numDate.Compare(time.Now()) < 0 {
		return nil, errors.New("TOKEN EXPIRED")
	}

	suuid, err := uuid.Parse(claims.SessionID)
	if err != nil {
		return nil, err
	}

	player, err := s.Repository.GetPlayerBySession(ctx, suuid)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func (s *Service) generateJWT(playerSession *models.PlayerSession) (string, error) {
	claims := ClaimType{
		SessionID: playerSession.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "game-integration-api-demo",
			ExpiresAt: jwt.NewNumericDate(playerSession.ExpiresAt),
			IssuedAt:  jwt.NewNumericDate(playerSession.IssuedAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := internal.Config.JWT_SECRET
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) validateJWT(tokenString string) (*ClaimType, error) {
	secretKey := internal.Config.JWT_SECRET

	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*ClaimType); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// validatePassword compares the provided password with the stored hash
func (s *Service) validatePassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// HashPassword creates a bcrypt hash of the password (utility function for creating users)
func (s *Service) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}
