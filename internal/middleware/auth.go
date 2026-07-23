package middleware

import (
	"net/http"
	"strings"

	"goecommerce/internal/auth"
	"goecommerce/internal/httpresponse"

	"github.com/labstack/echo/v5"
)

func AuthMiddleware(jwtService auth.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			var tokenString string

			// 1. Try to get token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenString = parts[1]
				}
			}

			// 2. If not found in header, try to get from cookie
			if tokenString == "" {
				if cookie, err := c.Cookie("access_token"); err == nil && cookie != nil {
					tokenString = cookie.Value
				} else if cookie, err := c.Cookie("token"); err == nil && cookie != nil {
					tokenString = cookie.Value
				}
			}

			// 3. If token is still empty, return Unauthorized
			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, httpresponse.Error{
					Code:    "UNAUTHORIZED",
					Message: "Missing or invalid authorization token",
				})
			}
			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, httpresponse.Error{
					Code:    "UNAUTHORIZED",
					Message: "Invalid or expired token",
				})
			}

			if claims.TokenType != auth.TokenTypeAccess {
				return c.JSON(http.StatusUnauthorized, httpresponse.Error{
					Code:    "UNAUTHORIZED",
					Message: "Invalid token type, expected access token",
				})
			}

			// Store the user id in context for downstream handlers to use
			c.Set("user_id", claims.UserId)

			return next(c)
		}
	}
}
