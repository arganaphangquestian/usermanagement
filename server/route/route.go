package route

import (
	"fmt"
	"strings"

	"github.com/arganaphangquestian/usermanagement/server/model"
	"github.com/arganaphangquestian/usermanagement/server/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type userService struct {
	service service.UserService
}

func (s *userService) register(c *fiber.Ctx) error {
	p := new(model.InputUser)
	if err := c.BodyParser(p); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}
	_, err := s.service.Register(*p)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}
	token, err := s.service.Login(model.Login{Username: p.Username, Password: p.Password})
	// DO ANOTHER PROCESS
	s.service.CreateSubcriber(token.AccessToken)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}
	cookie := new(fiber.Cookie)
	cookie.Name = "jid"
	cookie.Value = token.RefreshToken
	cookie.HTTPOnly = true
	c.Cookie(cookie)
	return c.Status(201).JSON(&fiber.Map{
		"success": true,
		"message": "Login Successfully",
		"data":    &token,
	})
}

func (s *userService) users(c *fiber.Ctx) error {
	response, err := s.service.Users()
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}
	return c.Status(201).JSON(&fiber.Map{
		"success": true,
		"message": "Get All Users endpoint reached",
		"data":    response,
	})
}

func (s *userService) login(c *fiber.Ctx) error {
	p := new(model.Login)
	if err := c.BodyParser(p); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}
	token, err := s.service.Login(*p)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}
	cookie := new(fiber.Cookie)
	cookie.Name = "jid"
	cookie.Value = token.RefreshToken
	cookie.HTTPOnly = true
	c.Cookie(cookie)
	return c.Status(201).JSON(&fiber.Map{
		"success": true,
		"message": "Login Successfully",
		"data":    &token,
	})
}

func (s *userService) dashboard(c *fiber.Ctx) error {
	authorizationHeader := c.Get("Authorization")
	if !strings.Contains(authorizationHeader, "Bearer") {
		return c.Status(403).JSON(&fiber.Map{
			"success": false,
			"message": "Authorization must be valid",
		})
	}
	user, err := extractToken(authorizationHeader)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}
	return c.Status(201).JSON(&fiber.Map{
		"success": true,
		"message": "DASHBOARD",
		"data":    user,
	})
}

func extractToken(authorizationHeader string) (interface{}, error) {
	claims := jwt.MapClaims{}
	tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Signing method invalid")
		} else if method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("Signing method invalid")
		}
		return []byte("MY_SUPER_SECRET_KEY"), nil
	})
	if err != nil {
		return nil, fmt.Errorf("Error : %s", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("Token not valid")
	}
	return claims["user"], nil
}

// New Route
func New(service service.UserService) *fiber.App {
	app := fiber.New()
	repo := &userService{service}
	app.Get("/user", repo.users)
	app.Post("/register", repo.register)
	app.Post("/login", repo.login)
	app.Get("/dashboard", repo.dashboard)
	return app
}
