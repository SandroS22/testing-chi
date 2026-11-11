package service

import (
	"time"

	customerrors "github.com/SandroS22/chi-test/cmd/CustomErrors"
	"github.com/SandroS22/chi-test/cmd/request"
	"github.com/SandroS22/chi-test/cmd/response"
	"github.com/SandroS22/chi-test/db"
	"github.com/SandroS22/chi-test/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("really_secret_key_with_a_lot_of_characters_981723813") // Load this from a secure location in a real application
var expirationTime = 24 * 60 * 60                                              // 24 hours in seconds

type AuthService interface {
	Login(req request.LoginRequest) (response.LoginResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (a *authService) Login(req request.LoginRequest) (response.LoginResponse, error) {
	user, err := a.userRepo.GetByUsername(req.Username)
	if err != nil {
		return response.LoginResponse{}, customerrors.ErrInvalidCredentials
	}

	if !CheckPasswordHash(req.Password, user.Password, user) {
		return response.LoginResponse{}, customerrors.ErrInvalidCredentials
	}

	jwtToken, err := GenerateJWT(user.ID, user.Role)
	if err != nil {
		return response.LoginResponse{}, err
	}

	return response.LoginResponse{
		Token: jwtToken,
	}, nil
}

func CheckPasswordHash(s1, s2 string, user *db.User) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(s1))
	if err != nil {
		return false
	}
	return true
}

func GenerateJWT(userID int, userRole string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    userRole,
		"exp":     time.Now().Add(time.Duration(expirationTime) * time.Second).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
