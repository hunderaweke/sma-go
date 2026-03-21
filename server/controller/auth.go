package controller

import (
	"fmt"
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

type AuthController struct {
	usecase domain.UserUsecase
}

func (c *AuthController) configureAuth() {
	goth.UseProviders(
		github.New(
			config.GitHubClientID,
			config.GitHubClientSecret,
			config.GitHubCallbackURL,
			"user:email",
			"user:profile",
		),
		google.New(
			config.GoogleClientID,
			config.GoogleClientSecret,
			config.GoogleCallbackURL,
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

func NewAuthController(uc domain.UserUsecase) *AuthController {
	ctrl := &AuthController{usecase: uc}
	ctrl.configureAuth()
	return ctrl
}

func (uc *AuthController) SignUpOrLogIn(c *fiber.Ctx) error {
	return goth_fiber.BeginAuthHandler(c)
}

func (uc *AuthController) AuthCallback(c *fiber.Ctx) error {
	user, err := goth_fiber.CompleteUserAuth(c)
	if err != nil {
		return c.Redirect(fmt.Sprintf("%s?error=authentication_failed", config.WebUrl), fiber.StatusExpectationFailed)
	}
	dbUser, err := uc.usecase.GetByEmail(user.Email)
	if dbUser != nil && (dbUser.Provider != user.Provider || dbUser.ProviderUserID != user.UserID) {
		return c.Redirect(fmt.Sprintf("%s?error=email_registered_with_different_provider", config.WebUrl), fiber.StatusConflict)
	}
	if dbUser == nil {
		dbUser, err = uc.usecase.Create(domain.User{
			Name:           user.Name,
			Provider:       user.Provider,
			ProviderUserID: user.UserID,
			Email:          user.Email,
		})
		if err != nil {
			return writeDomainError(c, err)
		}
	}
	accessToken, err := services.CreateAccessToken(dbUser.ID.String(), dbUser.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create access token"})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Expires:  time.Now().Add(15 * time.Hour),
	})
	return c.Redirect(config.WebUrl, fiber.StatusFound)
}

func (uc *AuthController) GetMe(c *fiber.Ctx) error {
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
func (uc *AuthController) Logout(c *fiber.Ctx) error {
	session, err := goth_fiber.SessionStore.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get session"})
	}
	if err := session.Destroy(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to destroy session"})
	}
	c.ClearCookie("session_id")
	c.ClearCookie("access_token")
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Expires:  time.Now().Add(-time.Hour),
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logged out successfully"})
}
