package middlewares

import (
	"spotsync/internal/apperror"
	"spotsync/internal/auth"
	"spotsync/internal/ctxkeys"
	"strings"

	"github.com/labstack/echo/v5"
)

func AuthMiddleware(jwtService auth.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {

			// extract token from authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return apperror.NewUnauthorized(nil, "Missing authorization header")
			}

			// check bearer scheme
			parts := strings.Split(authHeader, " ")

			if len(parts) != 2 || parts[0] != "Bearer" {
				return apperror.NewUnauthorized(nil, "invalid authorization header format")
			}

			tokenString := parts[1]

			// validate token

			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				return apperror.NewUnauthorized(err, "invalid or expired token")
			}

			// store user info in context for handlers
			c.Set(string(ctxkeys.UserID), claims.UserID)
			c.Set(string(ctxkeys.UserRole), claims.Role)

			return next(c)
		}
	}
}
