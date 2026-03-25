// Package controllers implements the HTTP handlers for all API endpoints.
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
	UserID   uint   `json:"UserID"`
	RoleName string `json:"RoleName"`
	jwt.RegisteredClaims
}

// Login authenticates a user by email and password, then returns a signed JWT token.
// The token contains the user's ID and role name, and expires after 2 hours.
// Deactivated users are rejected even if credentials are valid.
// Both email-not-found and wrong-password return the same error to prevent user enumeration.
//
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

	// Bind login credentials
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Check for email
	var existingUser models.Users
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	// Check for Password
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	// Check if user account is active
	if !existingUser.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account is deactivated"})
		return
	}

	// Fetch the user's role to embed in the token
	var role models.Roles
	if err := config.DB.First(&role, existingUser.RolesID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user role"})
		return
	}

	// Build JWT claims with user ID and role
	claim := &CustomClaim{
		UserID:   existingUser.ID,
		RoleName: role.RoleName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, tokenString)

}

// CreateUser registers a new staff user.
// Validates email uniqueness, password strength (via utils.ValidatePassword), and that the
// referenced role exists. The password is bcrypt-hashed before storage.
//
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Check if email already in db
	var count int64
	config.DB.Model(&models.Users{}).Where("email = ?", input.Email).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already in use"})
		return
	}

	if err := utils.ValidatePassword(input.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate and save token from password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User Created"})
}

// DeleteUser soft-deletes a user (sets deleted_at timestamp).
// The user record is preserved for order audit trails but hidden from all queries.
// The user's email becomes available for reuse after deletion.
//
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Check if the user exists
	if err := config.DB.Preload("Role").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Prevent deleting the last active admin
	if user.Role.RoleName == "admin" && user.IsActive {
		var adminCount int64
		config.DB.Model(&models.Users{}).
			Joins("JOIN roles ON roles.id = users.roles_id").
			Where("roles.role_name = ? AND users.is_active = ? AND users.deleted_at IS NULL", "admin", true).
			Count(&adminCount)
		if adminCount <= 1 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete the last active admin account"})
			return
		}
	}

	// Soft-delete user (sets deleted_at, preserves record for order audit trails)
	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

// ToggleUserStatus flips the IsActive flag on a user account.
// Deactivated users cannot log in but their data and order history are preserved.
// This is the preferred way to revoke access without losing audit trails.
//
// @Summary Toggle user active status
// @Description Activate or deactivate a user account without deleting it
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.Users
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /users/{id}/status [patch]
func ToggleUserStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var user models.Users
	if err := config.DB.Preload("Role").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Prevent deactivating the last active admin
	if user.Role.RoleName == "admin" && user.IsActive {
		var adminCount int64
		config.DB.Model(&models.Users{}).
			Joins("JOIN roles ON roles.id = users.roles_id").
			Where("roles.role_name = ? AND users.is_active = ? AND users.deleted_at IS NULL", "admin", true).
			Count(&adminCount)
		if adminCount <= 1 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot deactivate the last active admin account"})
			return
		}
	}

	user.IsActive = !user.IsActive

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user status"})
		return
	}

	config.DB.Preload("Role").First(&user, id)

	c.JSON(http.StatusOK, user)
}

// ChangePassword updates a user's password after verifying the current one.
// Non-admin users can only change their own password. Admins can change any user's password.
// The new password must pass the same strength validation as during user creation.
//
// @Summary Change user password
// @Description Change a user's password by providing the current and new password
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param passwords body object true "Old and new password"
// @Success 200 {object} map[string]string "Password updated"
// @Failure 400 {object} map[string]string "Invalid data"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "User not found"
// @Security BearerAuth
// @Router /users/{id}/password [patch]
func ChangePassword(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Non-admin users can only change their own password
	currentUserID := c.GetInt("userID")
	role := c.GetString("userRole")
	if role != "admin" && currentUserID != id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only change your own password"})
		return
	}

	var input struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	var user models.Users
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Current password is incorrect"})
		return
	}

	// Validate new password
	if err := utils.ValidatePassword(input.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash and save
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	if err := config.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated"})
}

// ResetPassword allows an admin to reset any user's password to a temporary value.
// The user should change this password on their next login via the ChangePassword endpoint.
//
// @Summary Reset user password (admin only)
// @Description Admin resets a user's password and receives a temporary password
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string "Temporary password"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Server error"
// @Security BearerAuth
// @Router /users/{id}/reset-password [patch]
func ResetPassword(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var user models.Users
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Generate a random temporary password
	tempPassword, err := utils.GenerateTempPassword(12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate temporary password"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	if err := config.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Password reset successful",
		"temp_password":  tempPassword,
	})
}

// GetUsers returns all users with their associated role preloaded.
// Password fields are excluded from the JSON response via the json:"-" tag on the model.
//
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUser returns a single user by ID with their role preloaded.
//
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := config.DB.Preload("Role").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, user)
}
