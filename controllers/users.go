package controllers

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"
	"wacdo/config"
	"wacdo/models"
	"wacdo/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CustomClaim struct {
	UserID uint
	jwt.RegisteredClaims
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Users
// @Accept json
// @Produce json
// @Param credentials body models.Users true "User credentials (email and password)"
// @Success 200 {string} string "JWT token"
// @Failure 400 {object} map[string]string "Invalid credentials"
// @Router /users/login [post]
func Login(c *gin.Context) {

	// Create the user object
	var user models.Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Data"})
		return
	}

	// Check for email
	var existingUser models.Users
	if err := config.DB.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or Password Invalid"})
		return
	}

	// Check for Password
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or Password Invalid"})
		return
	}

	//if both email ans Password are correct
	claim := &CustomClaim{
		UserID: existingUser.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while generating the token "})
		return
	}

	c.JSON(http.StatusOK, tokenString)

}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.UserInput true "User details"
// @Success 200 {object} map[string]string "User Created"
// @Failure 400 {object} map[string]string "Invalid data or email already in use"
// @Failure 500 {object} map[string]string "Server error"
// @Router /users [post]
func CreateUser(c *gin.Context) {

	// Get input
	var input models.UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Data"})
		return
	}

	// Check if email already in db
	var count int64
	config.DB.Model(&models.Users{}).Where("email = ?", input.Email).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email Already in Use"})
		return
	}

	if err := utils.ValidatePassword(input.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate and save token from password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server Error"})
		return
	}

	// Check role exists
	var role models.Roles
	if err := config.DB.First(&role, input.RolesID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role not found"})
		return
	}

	// Create the model from input
	user := models.Users{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		RolesID:  role.ID,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User Creation Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User Created"})
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string "User deleted"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	var user models.Users

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	// Check if id format is valid
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	// Check if the user exists
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User Not found"})
		return
	}

	// remove user
	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't delete user, internal error!"})
		return
	}
}

// GetUsers godoc
// @Summary Get all users
// @Description Retrieve a list of all users with their roles
// @Tags Users
// @Produce json
// @Success 200 {array} models.Users
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /users [get]
func GetUsers(c *gin.Context) {

	var users []models.Users

	if err := config.DB.Preload("Role").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't get User data"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Retrieve a single user by their ID
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.Users
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {

	var user models.Users

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := config.DB.Preload("Role").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User can't be found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, user)
}
