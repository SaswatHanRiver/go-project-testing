package services

import (
	"errors"
	"log/slog"

	"go-project-testing/models"
	"go-project-testing/repositories"
	"go-project-testing/utils"

	"golang.org/x/crypto/bcrypt"
)

// AuthService - equivalent to UserDetailsService + AuthenticationManager in Spring Boot
type AuthService struct {
	repo repositories.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{repo: repositories.UserRepository{}}
}

// Register - creates a new user account
// Equivalent to authService.register() in a Spring Boot AuthService
func (s *AuthService) Register(req *models.RegisterRequest) (models.AuthResponse, error) {
	// Check duplicate email - like Spring's existsByEmail()
	if s.repo.ExistsByEmail(req.Email) {
		return models.AuthResponse{}, errors.New("email already registered")
	}

	// Hash password with bcrypt - exactly like BCryptPasswordEncoder.encode() in Spring Boot
	// cost=14 means 2^14 hashing rounds (Spring Boot default is 10)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		return models.AuthResponse{}, errors.New("failed to hash password")
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(&user); err != nil {
		return models.AuthResponse{}, errors.New("failed to create user")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return models.AuthResponse{}, errors.New("failed to generate token")
	}

	slog.Info("User registered", "email", user.Email, "userID", user.ID)

	return models.AuthResponse{Token: token, User: user}, nil
}

// Login - authenticates user and returns JWT token
// Equivalent to authenticationManager.authenticate() + jwtUtil.generateToken() in Spring Boot
func (s *AuthService) Login(req *models.LoginRequest) (models.AuthResponse, error) {
	// Find user by email
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		// Return generic message - don't reveal whether email exists
		return models.AuthResponse{}, errors.New("invalid email or password")
	}

	// Compare hashed password - like BCryptPasswordEncoder.matches() in Spring Boot
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return models.AuthResponse{}, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return models.AuthResponse{}, errors.New("failed to generate token")
	}

	slog.Info("User logged in", "email", user.Email, "userID", user.ID)

	return models.AuthResponse{Token: token, User: user}, nil
}
