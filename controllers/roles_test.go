package controllers

import (
	"net/http"
	"testing"
	"wacdo/testutils"

	"github.com/stretchr/testify/assert"
)

func TestCreateRole_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	_ = db

	r := testutils.SetupRouter()
	r.POST("/roles", CreateRole)

	body := map[string]string{"role_name": "admin", "description": "Admin role"}
	req := testutils.JSONRequest("POST", "/roles", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := testutils.ParseResponse(w)
	assert.Equal(t, "admin", resp["role_name"])
}

func TestCreateRole_Duplicate(t *testing.T) {
	db := testutils.SetupTestDB()
	testutils.SeedRole(db, "admin")

	r := testutils.SetupRouter()
	r.POST("/roles", CreateRole)

	body := map[string]string{"role_name": "admin"}
	req := testutils.JSONRequest("POST", "/roles", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := testutils.ParseResponse(w)
	assert.Contains(t, resp["error"], "already exists")
}

func TestCreateRole_InvalidData(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.POST("/roles", CreateRole)

	req := testutils.JSONRequest("POST", "/roles", "invalid")
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetRoles_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	testutils.SeedRole(db, "admin")
	testutils.SeedRole(db, "accueil")

	r := testutils.SetupRouter()
	r.GET("/roles", GetRoles)

	req := testutils.JSONRequest("GET", "/roles", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetRole_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")

	r := testutils.SetupRouter()
	r.GET("/roles/:id", GetRole)

	req := testutils.JSONRequest("GET", testutils.IDParam("/roles", role.ID), nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := testutils.ParseResponse(w)
	assert.Equal(t, "admin", resp["role_name"])
}

func TestGetRole_NotFound(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.GET("/roles/:id", GetRole)

	req := testutils.JSONRequest("GET", "/roles/999", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetRole_InvalidID(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.GET("/roles/:id", GetRole)

	req := testutils.JSONRequest("GET", "/roles/abc", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteRole_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "temp")

	r := testutils.SetupRouter()
	r.DELETE("/roles/:id", DeleteRole)

	req := testutils.JSONRequest("DELETE", testutils.IDParam("/roles", role.ID), nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteRole_NotFound(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.DELETE("/roles/:id", DeleteRole)

	req := testutils.JSONRequest("DELETE", "/roles/999", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteRole_StillInUse(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	testutils.SeedUser(db, "user1", "user@test.com", "P@ssw0rd", role.ID)

	r := testutils.SetupRouter()
	r.DELETE("/roles/:id", DeleteRole)

	req := testutils.JSONRequest("DELETE", testutils.IDParam("/roles", role.ID), nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	resp := testutils.ParseResponse(w)
	assert.Contains(t, resp["error"], "still in use")
}
