package controller

import (
	"log"
	"strings"
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
			"read:user",
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
		CookieSecure:   config.Env == config.Production,
		CookieHTTPOnly: true,
		CookieSameSite: fiber.CookieSameSiteNoneMode,
		CookiePath:     "/",
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
	log.Printf("oauth start method=%s path=%s provider=%s host=%s referer=%s user_agent=%s session_id_present=%t session_id_len=%d",
		c.Method(),
		c.OriginalURL(),
		c.Params("provider"),
		c.Hostname(),
		c.Get("Referer"),
		c.Get("User-Agent"),
		c.Cookies("session_id") != "",
		len(c.Cookies("session_id")),
	)
	return goth_fiber.BeginAuthHandler(c)
}

func (uc *AuthController) AuthCallback(c *fiber.Ctx) error {
	log.Printf("oauth callback method=%s path=%s provider=%s host=%s referer=%s user_agent=%s session_id_present=%t session_id_len=%d query_error=%s query_code=%s",
		c.Method(),
		c.OriginalURL(),
		c.Params("provider"),
		c.Hostname(),
		c.Get("Referer"),
		c.Get("User-Agent"),
		c.Cookies("session_id") != "",
		len(c.Cookies("session_id")),
		c.Query("error"),
		c.Query("code"),
	)
	user, err := goth_fiber.CompleteUserAuth(c)
	if err != nil {
		log.Printf("authentication failed provider=%s path=%s error=%v session_id_present=%t session_id_len=%d",
			c.Params("provider"),
			c.OriginalURL(),
			err,
			c.Cookies("session_id") != "",
			len(c.Cookies("session_id")),
		)
		return c.Redirect(frontendURL("?error=authentication_failed"), fiber.StatusFound)
	}
	log.Printf("oauth user resolved provider=%s email=%s user_id=%s name=%s", user.Provider, user.Email, user.UserID, user.Name)
	dbUser, err := uc.usecase.GetByEmail(user.Email)
	if dbUser != nil && (dbUser.Provider != user.Provider || dbUser.ProviderUserID != user.UserID) {
		log.Printf("email registered with different provider email=%s provider=%s existing_provider=%s existing_provider_user_id=%s", user.Email, user.Provider, dbUser.Provider, dbUser.ProviderUserID)
		return c.Redirect(frontendURL("?error=email_registered_with_different_provider"), fiber.StatusFound)
	}
	if dbUser == nil {
		name := user.Name
		if name == "" {
			name = strings.Split(user.Email, "@")[0]
		}
		log.Printf("creating user from oauth email=%s provider=%s provider_user_id=%s name=%s", user.Email, user.Provider, user.UserID, name)
		dbUser, err = uc.usecase.Create(domain.User{
			Name:           name,
			Provider:       user.Provider,
			ProviderUserID: user.UserID,
			Email:          user.Email,
		})
		if err != nil {
			log.Printf("failed creating oauth user email=%s provider=%s provider_user_id=%s error=%v", user.Email, user.Provider, user.UserID, err)
			return writeDomainError(c, err)
		}
	}
	accessToken, err := services.CreateAccessToken(dbUser.ID.String(), dbUser.Email)
	if err != nil {
		log.Printf("failed creating access token user_id=%s email=%s error=%v", dbUser.ID.String(), dbUser.Email, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create access token"})
	}
	log.Printf("oauth success user_id=%s email=%s redirect_to=%s", dbUser.ID.String(), dbUser.Email, frontendURL(""))
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Path:     "/",
		Value:    accessToken,
		HTTPOnly: true,
		Secure:   config.Env == config.Production,
		SameSite: fiber.CookieSameSiteNoneMode,
		Expires:  time.Now().Add(15 * time.Hour),
	})
	return c.Redirect(frontendURL(""), fiber.StatusFound)
}

func frontendURL(query string) string {
	base := strings.TrimRight(config.WebUrl, "/")
	return base + query
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
		Path:     "/",
		Secure:   config.Env == config.Production,
		SameSite: fiber.CookieSameSiteNoneMode,
		Expires:  time.Now().Add(-time.Hour),
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logged out successfully"})
}
