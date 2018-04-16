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
	"google.golang.org/grpc/metadata"
)

// Login handlers
func Login(c echo.Context) error {
	// Use the following 2 lines for form body
	// username := c.FormValue("username")
	// password := c.FormValue("password")

	// Use the following code piece for json body
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	username := m["username"].(string)
	password := m["password"].(string)

	cn := pb.NewUserClient(Connections["user"])

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Attach rid to ctx, so that the micro service can access the same rid
	ctx = metadata.AppendToOutgoingContext(ctx, "rid", c.Response().Header().Get(echo.HeaderXRequestID))
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
