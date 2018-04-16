package main

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/rod6/tron/gateway/conf"
	"github.com/rod6/tron/gateway/ctrl"
)

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	ctrl.Logger = e.Logger

	defer ctrl.CloseConn()
	// Connect to micro services
	if err := ctrl.ConnSrv(); err != nil {
		e.Logger.Fatal("exits due to service connect error: %v", err)
	}
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{Generator: func() string {
		id := uuid.New()
		return id.String()
	}}))

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "hello, world",
		})
	})

	// curl localhost:1323/login -X POST -H "Content-type: application/json" -d '{"username":"rod","password":"123456"}'
	e.POST("/login", ctrl.Login)

	// Restricted group
	r := e.Group("/auth")
	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Skipper:       middleware.DefaultSkipper,
		SigningMethod: middleware.AlgorithmHS256,
		ContextKey:    "token",
		TokenLookup:   "header:" + echo.HeaderAuthorization,
		AuthScheme:    "Bearer",
		Claims:        jwt.MapClaims{},
		SigningKey:    []byte(conf.SignKey),
	}))

	// curl localhost:1323/auth/test -H "Authorization: Bearer xxxxxx"
	r.GET("/test", ctrl.Auth)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
