package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"wacdo/config"
	"wacdo/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB creates an in-memory SQLite database and runs migrations.
// It sets config.DB so controllers work without changes.
func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}

	db.AutoMigrate(
		&models.Roles{},
		&models.Users{},
		&models.Category{},
		&models.Products{},
		&models.ProductOptions{},
		&models.OptionValues{},
		&models.Menu{},
		&models.MenuProduct{},
		&models.Customer{},
		&models.Order{},
		&models.OrderItem{},
		&models.OrderItemOption{},
	)

	config.DB = db

	// Set a test JWT secret
	os.Setenv("JWT_SECRET", "test-secret-key")

	return db
}

// SetupRouter creates a gin engine in test mode.
func SetupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// SeedRole creates a role in the test DB and returns it.
func SeedRole(db *gorm.DB, name string) models.Roles {
	role := models.Roles{RoleName: name, Description: name + " role"}
	db.Create(&role)
	return role
}

// SeedUser creates a user with a hashed password in the test DB and returns it.
func SeedUser(db *gorm.DB, username, email, password string, roleID uint) models.Users {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.Users{
		Username: username,
		Email:    email,
		Password: string(hashed),
		RolesID:  roleID,
		IsActive: true,
	}
	db.Create(&user)
	return user
}

// SeedCategory creates a category in the test DB.
func SeedCategory(db *gorm.DB, name string) models.Category {
	cat := models.Category{Name: name, Description: name + " category"}
	db.Create(&cat)
	return cat
}

// SeedProduct creates a product in the test DB.
func SeedProduct(db *gorm.DB, name string, price float64, categoryID uint, available bool) models.Products {
	p := models.Products{
		Name:        name,
		Price:       price,
		CategoryID:  categoryID,
		IsAvailable: available,
	}
	db.Create(&p)
	return p
}

// SeedMenu creates a menu in the test DB.
func SeedMenu(db *gorm.DB, name string, price float64, available bool) models.Menu {
	m := models.Menu{Name: name, Price: price, IsAvailable: available}
	db.Create(&m)
	return m
}

// SeedCustomer creates a customer in the test DB.
func SeedCustomer(db *gorm.DB, name, phone, email string) models.Customer {
	c := models.Customer{Name: name, Phone: phone, Email: email}
	db.Create(&c)
	return c
}

// JSONRequest builds an HTTP request with a JSON body.
func JSONRequest(method, url string, body interface{}) *http.Request {
	jsonBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(method, url, bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	return req
}

// PerformRequest executes a request against a gin router and returns the recorder.
func PerformRequest(r *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// ParseResponse unmarshals a JSON response body into a map.
func ParseResponse(w *httptest.ResponseRecorder) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	return result
}

// AuthMiddleware is a test middleware that sets userID and userRole in the gin context.
func AuthMiddleware(userID int, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("userRole", role)
		c.Next()
	}
}

// IDParam returns the URL with the id substituted (for use in route definitions).
func IDParam(base string, id uint) string {
	return fmt.Sprintf("%s/%d", base, id)
}

// GetDB returns the current test database instance.
func GetDB() *gorm.DB {
	return config.DB
}
