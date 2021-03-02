package service

import (
	"fmt"
	"os"
	"time"

	"github.com/arganaphangquestian/usermanagement/server/model"
	"github.com/arganaphangquestian/usermanagement/server/repository"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type (
	service struct {
		repo repository.UserRepository
	}
)

// Register methods
func (s *service) Register(register model.InputUser) (*model.User, error) {
	hashPassword, err := passwordHash(register.Password)
	if err != nil {
		return nil, err
	}
	register.Password = hashPassword
	return s.repo.Register(register)
}

// Users methods
func (s *service) Users() ([]*model.User, error) {
	return s.repo.Users()
}

// Login methods
func (s *service) Login(login model.Login) (*model.UserToken, error) {
	userRes, err := s.repo.GetUserByUsername(login.Username)
	if err != nil {
		return nil, err
	}

	passwordMatch := comparePassword(userRes.Password, login.Password)

	if !passwordMatch {
		return nil, fmt.Errorf("Incorrect Username or Password")
	}

	user := model.User{
		ID:       userRes.ID,
		Name:     userRes.Name,
		Username: userRes.Username,
		Email:    userRes.Email,
	}

	userToken, err := createToken(user.ID)
	if err != nil {
		return nil, err
	}

	return userToken, nil
}

func (s *service) CreateSubcriber(token string) error {
	fmt.Printf("Create Subcriber for token : %s", token)
	return nil
}

func (s *service) RefreshToken(refeshToken string) (*model.UserToken, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(refeshToken, claims, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error : %s", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("Token not valid")
	}
	newToken, err := createToken(claims["id"].(string))
	return newToken, err
}

func (s *service) Verify(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return fmt.Errorf("Error : %s", err)
	}
	if !token.Valid {
		return fmt.Errorf("Token not valid")
	}
	return nil
}

func passwordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func comparePassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func createToken(userID string) (*model.UserToken, error) {
	// Refresh Token
	refreshToken := jwt.MapClaims{}
	refreshToken["authorize"] = true
	refreshToken["id"] = userID
	refreshToken["exp"] = time.Now().Add(time.Hour * 24 * 30 * 12 * 5).Unix() // 5 years
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshToken)
	rToken, err := rt.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return nil, err
	}

	// Access Token
	accessToken := jwt.MapClaims{}
	accessToken["authorize"] = true
	accessToken["id"] = userID
	accessToken["exp"] = time.Now().Add(time.Minute * 5).Unix() // 5 Minutes
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessToken)
	aToken, err := at.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return nil, err
	}

	return &model.UserToken{
		RefreshToken: rToken,
		AccessToken:  aToken,
	}, nil
}

// New UserService
func New(repo repository.UserRepository) UserService {
	return &service{repo}
}
