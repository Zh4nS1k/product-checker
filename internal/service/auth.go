package service

import (
	"errors"
	"time"

	"barcode-checker/internal/model"
	"barcode-checker/internal/repository"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(username, email, password string) (*model.User, error)
	Login(email, password string) (string, error)
	ListUsers() ([]model.User, error)
	DeleteUser(id uint) error
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
	jwtExp    int
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string, jwtExp int) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExp:    jwtExp,
	}
}

func (s *authService) Register(username, email, password string) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * time.Duration(s.jwtExp)).Unix(),
	})

	return token.SignedString([]byte(s.jwtSecret))
}

func (s *authService) ListUsers() ([]model.User, error) {
	return s.userRepo.ListUsers()
}

func (s *authService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}
