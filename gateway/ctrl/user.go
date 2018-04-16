// Package ctrl for controllers used in gateway
package ctrl

import (
	"context"
	"errors"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/rod6/tron/gateway/conf"
	pb "github.com/rod6/tron/pb"
)

// Login handlers
func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	cn := pb.NewUserClient(Connections["user"])

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := cn.Auth(ctx, &pb.AuthRequest{Username: username, Password: password})
	if err != nil {
		Logger.Error("call micro service 'user' error: ", err)
		return errors.New("internal service error")
	}

	if r.Authed {
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["jti"] = username
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(conf.SignKey))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}

	return echo.ErrUnauthorized
}

// Auth handler
func Auth(c echo.Context) error {
	token := c.Get("token").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	name := claims["jti"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}
