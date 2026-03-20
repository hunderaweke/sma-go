package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/storage/sqlite3"
	"github.com/hunderaweke/sma-go/config"
	"github.com/hunderaweke/sma-go/domain"
	"github.com/hunderaweke/sma-go/services"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/shareed2k/goth_fiber"
)

type UserController struct {
	usecase domain.UserUsecase
}

func (c *UserController) configureAuth() {
	goth.UseProviders(
		github.New(
			config.GitHubClientID,
			config.GitHubClientSecret,
			"http://localhost:3000/auth/github/callback",
			"user:email",
			"user:profile",
		),
		google.New(
			config.GoogleClientID,
			config.GoogleClientSecret,
			"http://localhost:3000/auth/google/callback",
			"openid",
			"email",
			"profile",
		),
	)
	store := session.New(session.Config{
		Expiration:     24 * time.Hour,
		KeyLookup:      "cookie:session_id",
		CookieSecure:   false,
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
		KeyGenerator:   utils.UUIDv4,
		Storage:        sqlite3.New(sqlite3.Config{Database: "./session.db"}),
	})
	goth_fiber.SessionStore = store
}

func NewUserController(uc domain.UserUsecase) *UserController {
	ctrl := &UserController{usecase: uc}
	ctrl.configureAuth()
	return ctrl
}

func (uc *UserController) SignUpOrLogIn(c *fiber.Ctx) error {
	return goth_fiber.BeginAuthHandler(c)
}

func (uc *UserController) AuthCallback(c *fiber.Ctx) error {
	user, err := goth_fiber.CompleteUserAuth(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authentication failed"})
	}
	prevUser, err := uc.usecase.GetByEmail(user.Email)
	if err == nil {
		if prevUser.Provider != user.Provider || prevUser.ProviderUserID != user.UserID {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already registered with a different provider"})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User already exists, logged in successfully"})
	}
	dbUser, err := uc.usecase.Create(domain.User{
		Name:           user.Name,
		Provider:       user.Provider,
		ProviderUserID: user.UserID,
		Email:          user.Email,
	})
	if err != nil {
		return writeDomainError(c, err)
	}
	accessToken, err := services.CreateAccessToken(dbUser.ID.String(), dbUser.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create access token"})
	}
	refreshToken, err := services.CreateRefreshToken(dbUser.ID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create refresh token"})
	}
	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          dbUser,
	})
}

func (uc *UserController) GetMe(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	user, err := uc.usecase.GetById(userID)
	if err != nil {
		return writeDomainError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(user)
}
