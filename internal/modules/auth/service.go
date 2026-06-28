package auth

import (
	"context"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"restaurant-inventory-api/internal/config"
	apperrors "restaurant-inventory-api/internal/errors"
)

type Service struct {
	repo Repository
	cfg  config.AuthConfig
}

func NewService(repo Repository, cfg config.AuthConfig) Service {
	return Service{repo: repo, cfg: cfg}
}

func (s Service) Register(ctx context.Context, req RegisterRequest) (AuthResponse, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.FullName = strings.TrimSpace(req.FullName)

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResponse{}, err
	}

	user, err := s.repo.CreateUser(ctx, req, string(hash))
	if IsUniqueViolation(err) {
		return AuthResponse{}, apperrors.Conflict("email already exists", nil)
	}
	if IsForeignKeyViolation(err) {
		return AuthResponse{}, apperrors.BadRequest("restaurant or branch does not exist", nil)
	}
	if err != nil {
		return AuthResponse{}, err
	}

	return s.authResponse(user)
}

func (s Service) Login(ctx context.Context, req LoginRequest) (AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	user, passwordHash, err := s.repo.GetByEmailWithPassword(ctx, email)
	if IsNotFound(err) {
		return AuthResponse{}, apperrors.Unauthorized("invalid email or password")
	}
	if err != nil {
		return AuthResponse{}, err
	}
	if user.Status != "active" {
		return AuthResponse{}, apperrors.Forbidden("user is not active")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		return AuthResponse{}, apperrors.Unauthorized("invalid email or password")
	}

	return s.authResponse(user)
}

func (s Service) Me(ctx context.Context, userID string) (User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if IsNotFound(err) {
		return user, apperrors.NotFound("user not found")
	}
	return user, err
}

func (s Service) authResponse(user User) (AuthResponse, error) {
	expiresAt := time.Now().Add(s.cfg.AccessTokenExpiry)
	claims := jwt.MapClaims{
		"sub":           user.ID,
		"email":         user.Email,
		"roles":         user.Roles,
		"restaurant_id": user.RestaurantID,
		"branch_id":     user.BranchID,
		"exp":           expiresAt.Unix(),
		"iat":           time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{
		AccessToken: signed,
		TokenType:   "Bearer",
		ExpiresIn:   int64(s.cfg.AccessTokenExpiry.Seconds()),
		User:        user,
	}, nil
}

func ParseToken(tokenString string, secret string) (Claims, error) {
	parsed, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !parsed.Valid {
		return Claims{}, apperrors.Unauthorized("invalid token")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return Claims{}, apperrors.Unauthorized("invalid token claims")
	}

	result := Claims{}
	if sub, ok := claims["sub"].(string); ok {
		result.UserID = sub
	}
	if email, ok := claims["email"].(string); ok {
		result.Email = email
	}
	if restaurantID, ok := claims["restaurant_id"].(string); ok && restaurantID != "" {
		result.RestaurantID = &restaurantID
	}
	if branchID, ok := claims["branch_id"].(string); ok && branchID != "" {
		result.BranchID = &branchID
	}
	if roles, ok := claims["roles"].([]interface{}); ok {
		for _, role := range roles {
			if value, ok := role.(string); ok {
				result.Roles = append(result.Roles, value)
			}
		}
	}

	if result.UserID == "" {
		return Claims{}, apperrors.Unauthorized("invalid token subject")
	}

	return result, nil
}
