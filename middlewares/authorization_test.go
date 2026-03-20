package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupAuthzRouter(allowedRoles ...string) *gin.Engine {
	r := gin.New()
	// Simulate authentication by setting userRole in context
	r.Use(func(c *gin.Context) {
		if role := c.GetHeader("X-Test-Role"); role != "" {
			c.Set("userRole", role)
		}
		c.Next()
	})
	r.Use(Authorization(allowedRoles...))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	return r
}

func TestAuthorization_AllowedRole(t *testing.T) {
	r := setupAuthzRouter("admin", "accueil")

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Test-Role", "admin")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthorization_DeniedRole(t *testing.T) {
	r := setupAuthzRouter("admin")

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Test-Role", "preparation")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestAuthorization_NoRoleInContext(t *testing.T) {
	r := setupAuthzRouter("admin")

	req := httptest.NewRequest("GET", "/test", nil)
	// No X-Test-Role header, so no userRole in context
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestAuthorization_MultipleAllowedRoles(t *testing.T) {
	r := setupAuthzRouter("admin", "accueil", "preparation")

	roles := []string{"admin", "accueil", "preparation"}
	for _, role := range roles {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Test-Role", role)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "role %s should be allowed", role)
	}
}

func TestAuthorization_EmptyStringRole(t *testing.T) {
	r := setupAuthzRouter("admin")

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Test-Role", "")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
