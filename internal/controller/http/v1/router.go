package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/s4lat/gosavingsbot/internal/log"
	"github.com/s4lat/gosavingsbot/internal/usecase"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

//	@title			Savings API
//	@version		1.0
//	@description	Savings Web App backend
//
// @BasePath	/v1
func NewRouter(uc *usecase.UseCase, botToken string, withSwagger bool) http.Handler {
	e := echo.New()
	e.Use(LogMiddleware())
	e.Use(middleware.Recover())

	healthzRoutes := NewHealthzRoutes()
	userRoutes := NewUserRoutes(uc)

	// TEST
	e.GET("/", userRoutes.Index)

	// MW
	tgWebAppAuthMW := TgWebAppAuthMiddleware(botToken)
	tgWebAppUserToServiceUserMW := TgWebAppUserToServiceUserMiddleware(uc)

	// V1 GROUP
	v1Group := e.Group("/v1")

	// USER
	v1Group.GET("/users/me", userRoutes.GetMe, tgWebAppAuthMW, tgWebAppUserToServiceUserMW)
	v1Group.POST("/users", userRoutes.CreateNewUser, tgWebAppAuthMW)

	// HEALTHCHECK
	v1Group.GET("/healthz", healthzRoutes.getHealthz)

	if withSwagger {
		v1Group.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	for _, route := range e.Routes() {
		log.Sugar().Infof("%s - %s", route.Method, route.Path)
	}
	return e
}
