package middlewares

import (
	"spotsync/internal/apperror"
	"spotsync/internal/ctxkeys"

	"github.com/labstack/echo/v5"
)

// RequireRole returns a middleware that checks the user's role against allowed roles.
// Must be used AFTER AuthMiddleware so that UserRole is set in context.
func RequireRole(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			role, ok := c.Get(string(ctxkeys.UserRole)).(string)
			if !ok || role == "" {
				return apperror.NewForbidden(nil, "access denied: role not found")
			}

			for _, allowed := range allowedRoles {
				if role == allowed {
					return next(c)
				}
			}

			return apperror.NewForbidden(nil, "access denied: insufficient permissions")
		}
	}
}
