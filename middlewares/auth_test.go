package middlewares

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "test-secret-key")
}

func generateToken(userID float64, roleName string, expiry time.Duration) string {
	claims := jwt.MapClaims{
		"UserID":   userID,
		"RoleName": roleName,
		"exp":      jwt.NewNumericDate(time.Now().Add(expiry)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return tokenString
}

func setupAuthRouter() *gin.Engine {
	r := gin.New()
	r.Use(Authentication())
	r.GET("/test", func(c *gin.Context) {
		userID := c.GetInt("userID")
		role := c.GetString("userRole")
		c.JSON(http.StatusOK, gin.H{"userID": userID, "userRole": role})
	})
	return r
}

func TestAuthentication_ValidToken(t *testing.T) {
	r := setupAuthRouter()
	token := generateToken(1, "admin", 2*time.Hour)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthentication_MissingHeader(t *testing.T) {
	r := setupAuthRouter()

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthentication_NoBearerPrefix(t *testing.T) {
	r := setupAuthRouter()
	token := generateToken(1, "admin", 2*time.Hour)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", token) // missing "Bearer "
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthentication_ExpiredToken(t *testing.T) {
	r := setupAuthRouter()
	token := generateToken(1, "admin", -1*time.Hour) // expired

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthentication_InvalidToken(t *testing.T) {
	r := setupAuthRouter()

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-string")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthentication_WrongSigningMethod(t *testing.T) {
	r := setupAuthRouter()

	// Create a token signed with RSA instead of HMAC
	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"UserID":   float64(1),
		"RoleName": "admin",
		"exp":      jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
	})
	tokenString, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthentication_SetsContextValues(t *testing.T) {
	r := setupAuthRouter()
	token := generateToken(42, "preparation", 2*time.Hour)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"userID":42`)
	assert.Contains(t, w.Body.String(), `"userRole":"preparation"`)
}
