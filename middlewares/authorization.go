package middlewares

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

// Authorization checks that the authenticated user's role is in the allowed list.
// Must be used after Authentication() so that "userRole" is set in the context.
func Authorization(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		roleName, ok := role.(string)
		if !ok || !slices.Contains(allowedRoles, roleName) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this resource"})
			return
		}

		c.Next()
	}
}
