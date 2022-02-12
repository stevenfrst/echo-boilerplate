package middlewares

import (
	"echo-boilerplate/delivery"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"time"
)

type JwtCustomClaims struct {
	ID         int  `json:"id"`
	IsVerified bool `json:"is_verified"`
	jwt.StandardClaims
}

type ConfigJWT struct {
	SecretJWT       string
	ExpiresDuration int64
}

func (jwtConf *ConfigJWT) Init() middleware.JWTConfig {
	return middleware.JWTConfig{
		Claims:     &JwtCustomClaims{},
		SigningKey: []byte(jwtConf.SecretJWT),
		ErrorHandlerWithContext: middleware.JWTErrorHandlerWithContext(func(e error, c echo.Context) error {
			return delivery.ErrorResponse(c, http.StatusForbidden, "", e)
		}),
	}
}

func (jwtConf *ConfigJWT) GenerateToken(userID int, IsVerified bool) string {
	claims := JwtCustomClaims{
		userID,
		IsVerified,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(jwtConf.ExpiresDuration)).Unix(),
		},
	}
	initToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, _ := initToken.SignedString([]byte(jwtConf.SecretJWT))
	return token
}

func GetUser(c echo.Context) *JwtCustomClaims {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)
	return claims
}
