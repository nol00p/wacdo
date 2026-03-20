package controllers

import (
	"encoding/json"
	"net/http"
	"testing"
	"wacdo/testutils"

	"github.com/stretchr/testify/assert"
)

func TestLogin_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	testutils.SeedUser(db, "admin", "admin@test.com", "P@ssw0rd", role.ID)

	r := testutils.SetupRouter()
	r.POST("/users/login", Login)

	body := map[string]string{"email": "admin@test.com", "password": "P@ssw0rd"}
	req := testutils.JSONRequest("POST", "/users/login", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Response is a JWT string
	var token string
	json.Unmarshal(w.Body.Bytes(), &token)
	assert.NotEmpty(t, token)
}

func TestLogin_WrongPassword(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	testutils.SeedUser(db, "admin", "admin@test.com", "P@ssw0rd", role.ID)

	r := testutils.SetupRouter()
	r.POST("/users/login", Login)

	body := map[string]string{"email": "admin@test.com", "password": "WrongP@ss1"}
	req := testutils.JSONRequest("POST", "/users/login", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_EmailNotFound(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.POST("/users/login", Login)

	body := map[string]string{"email": "nobody@test.com", "password": "P@ssw0rd"}
	req := testutils.JSONRequest("POST", "/users/login", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_DeactivatedUser(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	user := testutils.SeedUser(db, "admin", "admin@test.com", "P@ssw0rd", role.ID)
	// Deactivate user
	db.Model(&user).Update("is_active", false)

	r := testutils.SetupRouter()
	r.POST("/users/login", Login)

	body := map[string]string{"email": "admin@test.com", "password": "P@ssw0rd"}
	req := testutils.JSONRequest("POST", "/users/login", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestLogin_InvalidData(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.POST("/users/login", Login)

	body := map[string]string{"email": "not-an-email"}
	req := testutils.JSONRequest("POST", "/users/login", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUser_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")

	r := testutils.SetupRouter()
	r.POST("/users", CreateUser)

	body := map[string]interface{}{
		"username": "newuser",
		"email":    "new@test.com",
		"password": "P@ssw0rd",
		"roles_id": role.ID,
	}
	req := testutils.JSONRequest("POST", "/users", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := testutils.ParseResponse(w)
	assert.Equal(t, "User Created", resp["message"])
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	testutils.SeedUser(db, "existing", "dup@test.com", "P@ssw0rd", role.ID)

	r := testutils.SetupRouter()
	r.POST("/users", CreateUser)

	body := map[string]interface{}{
		"username": "newuser",
		"email":    "dup@test.com",
		"password": "P@ssw0rd",
		"roles_id": role.ID,
	}
	req := testutils.JSONRequest("POST", "/users", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := testutils.ParseResponse(w)
	assert.Contains(t, resp["error"], "Email already in use")
}

func TestCreateUser_WeakPassword(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")

	r := testutils.SetupRouter()
	r.POST("/users", CreateUser)

	body := map[string]interface{}{
		"username": "newuser",
		"email":    "weak@test.com",
		"password": "weakpw",
		"roles_id": role.ID,
	}
	req := testutils.JSONRequest("POST", "/users", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUser_InvalidRole(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.POST("/users", CreateUser)

	body := map[string]interface{}{
		"username": "newuser",
		"email":    "new@test.com",
		"password": "P@ssw0rd",
		"roles_id": 999,
	}
	req := testutils.JSONRequest("POST", "/users", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := testutils.ParseResponse(w)
	assert.Contains(t, resp["error"], "Role not found")
}

func TestGetUsers_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	testutils.SeedUser(db, "u1", "u1@test.com", "P@ssw0rd", role.ID)
	testutils.SeedUser(db, "u2", "u2@test.com", "P@ssw0rd", role.ID)

	r := testutils.SetupRouter()
	r.GET("/users", GetUsers)

	req := testutils.JSONRequest("GET", "/users", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUser_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	user := testutils.SeedUser(db, "admin", "admin@test.com", "P@ssw0rd", role.ID)

	r := testutils.SetupRouter()
	r.GET("/users/:id", GetUser)

	req := testutils.JSONRequest("GET", testutils.IDParam("/users", user.ID), nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := testutils.ParseResponse(w)
	assert.Equal(t, "admin@test.com", resp["email"])
}

func TestGetUser_NotFound(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.GET("/users/:id", GetUser)

	req := testutils.JSONRequest("GET", "/users/999", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetUser_InvalidID(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.GET("/users/:id", GetUser)

	req := testutils.JSONRequest("GET", "/users/abc", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteUser_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	user := testutils.SeedUser(db, "todelete", "del@test.com", "P@ssw0rd", role.ID)

	r := testutils.SetupRouter()
	r.DELETE("/users/:id", DeleteUser)

	req := testutils.JSONRequest("DELETE", testutils.IDParam("/users", user.ID), nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteUser_NotFound(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.DELETE("/users/:id", DeleteUser)

	req := testutils.JSONRequest("DELETE", "/users/999", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestToggleUserStatus_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	user := testutils.SeedUser(db, "admin", "admin@test.com", "P@ssw0rd", role.ID)

	r := testutils.SetupRouter()
	r.PATCH("/users/:id/status", ToggleUserStatus)

	// User starts active, toggle should deactivate
	req := testutils.JSONRequest("PATCH", testutils.IDParam("/users", user.ID)+"/status", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := testutils.ParseResponse(w)
	assert.Equal(t, false, resp["is_active"])

	// Toggle again to reactivate
	req = testutils.JSONRequest("PATCH", testutils.IDParam("/users", user.ID)+"/status", nil)
	w = testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp = testutils.ParseResponse(w)
	assert.Equal(t, true, resp["is_active"])
}

func TestToggleUserStatus_NotFound(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.PATCH("/users/:id/status", ToggleUserStatus)

	req := testutils.JSONRequest("PATCH", "/users/999/status", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestChangePassword_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	user := testutils.SeedUser(db, "admin", "admin@test.com", "P@ssw0rd", role.ID)

	r := testutils.SetupRouter()
	r.Use(testutils.AuthMiddleware(int(user.ID), "admin"))
	r.PATCH("/users/:id/password", ChangePassword)

	body := map[string]string{
		"old_password": "P@ssw0rd",
		"new_password": "N3wP@ssw0rd",
	}
	req := testutils.JSONRequest("PATCH", testutils.IDParam("/users", user.ID)+"/password", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := testutils.ParseResponse(w)
	assert.Equal(t, "Password updated", resp["message"])
}

func TestChangePassword_WrongOldPassword(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	user := testutils.SeedUser(db, "admin", "admin@test.com", "P@ssw0rd", role.ID)

	r := testutils.SetupRouter()
	r.Use(testutils.AuthMiddleware(int(user.ID), "admin"))
	r.PATCH("/users/:id/password", ChangePassword)

	body := map[string]string{
		"old_password": "WrongP@ss1",
		"new_password": "N3wP@ssw0rd",
	}
	req := testutils.JSONRequest("PATCH", testutils.IDParam("/users", user.ID)+"/password", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChangePassword_OtherUserForbidden(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "accueil")
	user := testutils.SeedUser(db, "user1", "user1@test.com", "P@ssw0rd", role.ID)
	otherUser := testutils.SeedUser(db, "user2", "user2@test.com", "P@ssw0rd", role.ID)

	r := testutils.SetupRouter()
	// Authenticated as user1 (non-admin), trying to change user2's password
	r.Use(testutils.AuthMiddleware(int(user.ID), "accueil"))
	r.PATCH("/users/:id/password", ChangePassword)

	body := map[string]string{
		"old_password": "P@ssw0rd",
		"new_password": "N3wP@ssw0rd",
	}
	req := testutils.JSONRequest("PATCH", testutils.IDParam("/users", otherUser.ID)+"/password", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestChangePassword_AdminCanChangeOthers(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	admin := testutils.SeedUser(db, "admin", "admin@test.com", "P@ssw0rd", role.ID)
	otherUser := testutils.SeedUser(db, "user2", "user2@test.com", "P@ssw0rd", role.ID)

	r := testutils.SetupRouter()
	r.Use(testutils.AuthMiddleware(int(admin.ID), "admin"))
	r.PATCH("/users/:id/password", ChangePassword)

	body := map[string]string{
		"old_password": "P@ssw0rd",
		"new_password": "N3wP@ssw0rd",
	}
	req := testutils.JSONRequest("PATCH", testutils.IDParam("/users", otherUser.ID)+"/password", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
