package middlewares

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/moody/config"
	"github.com/moody/helpers"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil {
			return echo.ErrUnauthorized
		}

		tokenString := cookie.Value
		claims := &config.JWTClaim{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return config.JWT_KEY, nil
		})

		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			fmt.Println(v)
			switch v.Errors {
			case jwt.ValidationErrorSignatureInvalid:
				return echo.ErrUnauthorized
			case jwt.ValidationErrorExpired:
				// token expired
				response := map[string]interface{}{
					"code":    401,
					"message": "Unauthorized, Token expired!",
				}
				return helpers.Response(c, http.StatusUnauthorized, response)
			default:
				return echo.ErrUnauthorized
			}
		}

		if !token.Valid {
			return echo.ErrUnauthorized
		}
		c.Set("jwt_user_id", claims.ID)
		c.Set("jwt_username", claims.UserName)
		c.Set("jwt_email", claims.Email)

		return next(c)
	}
}
