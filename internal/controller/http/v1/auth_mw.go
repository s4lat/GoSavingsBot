package v1

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/s4lat/gosavingsbot/internal/domain"
	"github.com/s4lat/gosavingsbot/internal/usecase"
	"github.com/s4lat/gosavingsbot/pkg/webAppAuth"
	"net/http"
)

func TgWebAppAuthMiddleware(botToken string) echo.MiddlewareFunc {
	userAuthorizer := webAppAuth.NewAuthorizer(botToken)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			l := getLoggerFromEchoContext(c)

			authTokenEnc := c.Request().Header.Get("Authorization")
			if authTokenEnc == "" {
				return c.JSON(http.StatusUnauthorized, newErrorResponse("empty auth token"))
			}

			authToken, err := decodeBase64(authTokenEnc)
			if err != nil {
				l.Errorf("can't decode auth token, %s: %s", err, authToken)
				return c.JSON(http.StatusUnauthorized, newErrorResponse("invalid auth token"))
			}

			tgWebAppUser, err := userAuthorizer.AuthorizeUser(authToken)
			if err != nil {
				l.Errorf("can't authorize user, %s: %s", err, authToken)
				return c.JSON(http.StatusUnauthorized, newErrorResponse("invalid auth token"))
			}

			c.Set("tgWebAppUser", tgWebAppUser)

			return next(c)
		}
	}
}

func TgWebAppUserToServiceUserMiddleware(uc *usecase.UseCase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			l := getLoggerFromEchoContext(c)
			tgWebAppUser := c.Get("tgWebAppUser").(webAppAuth.TelegramWebAppUser)

			serviceUser, err := uc.User.GetUserById(c.Request().Context(), tgWebAppUser.Id)
			if err != nil {
				if errors.Is(err, domain.ErrNotFound) {
					return c.JSON(http.StatusUnauthorized, newErrorResponse("user not found"))
				}

				l.Errorf("can't create/get user: %s", err)
				return c.JSON(http.StatusInternalServerError, newErrorResponse("can't get user"))
			}

			c.Set("user", serviceUser)

			return next(c)
		}
	}
}

// decodeBase64 decodes a string from Base64 and then from URL encoding
func decodeBase64(encoded string) (string, error) {
	// Step 1: Decode from Base64
	base64Decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("error decoding Base64: %w", err)
	}

	return string(base64Decoded), nil
}
