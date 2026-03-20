package controllers

import (
	"net/http"
	"testing"
	"wacdo/testutils"

	"github.com/stretchr/testify/assert"
)

func TestCreateCategory_Success(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.POST("/categories", CreateCategory)

	body := map[string]string{"name": "Burgers", "description": "Burger category"}
	req := testutils.JSONRequest("POST", "/categories", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := testutils.ParseResponse(w)
	assert.Equal(t, "Burgers", resp["name"])
}

func TestCreateCategory_Duplicate(t *testing.T) {
	db := testutils.SetupTestDB()
	testutils.SeedCategory(db, "Burgers")

	r := testutils.SetupRouter()
	r.POST("/categories", CreateCategory)

	body := map[string]string{"name": "Burgers"}
	req := testutils.JSONRequest("POST", "/categories", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetCategories_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	testutils.SeedCategory(db, "Burgers")
	testutils.SeedCategory(db, "Drinks")

	r := testutils.SetupRouter()
	r.GET("/categories", GetCategories)

	req := testutils.JSONRequest("GET", "/categories", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetCategory_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	cat := testutils.SeedCategory(db, "Burgers")

	r := testutils.SetupRouter()
	r.GET("/categories/:id", GetCategory)

	req := testutils.JSONRequest("GET", testutils.IDParam("/categories", cat.ID), nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := testutils.ParseResponse(w)
	assert.Equal(t, "Burgers", resp["name"])
}

func TestGetCategory_NotFound(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.GET("/categories/:id", GetCategory)

	req := testutils.JSONRequest("GET", "/categories/999", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateCategory_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	cat := testutils.SeedCategory(db, "Burgers")

	r := testutils.SetupRouter()
	r.PUT("/categories/:id", UpdateCategory)

	body := map[string]string{"name": "Premium Burgers"}
	req := testutils.JSONRequest("PUT", testutils.IDParam("/categories", cat.ID), body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateCategory_NameConflict(t *testing.T) {
	db := testutils.SetupTestDB()
	testutils.SeedCategory(db, "Drinks")
	cat := testutils.SeedCategory(db, "Burgers")

	r := testutils.SetupRouter()
	r.PUT("/categories/:id", UpdateCategory)

	body := map[string]string{"name": "Drinks"}
	req := testutils.JSONRequest("PUT", testutils.IDParam("/categories", cat.ID), body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestDeleteCategory_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	cat := testutils.SeedCategory(db, "Burgers")

	r := testutils.SetupRouter()
	r.DELETE("/categories/:id", DeleteCategory)

	req := testutils.JSONRequest("DELETE", testutils.IDParam("/categories", cat.ID), nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteCategory_StillInUse(t *testing.T) {
	db := testutils.SetupTestDB()
	cat := testutils.SeedCategory(db, "Burgers")
	testutils.SeedProduct(db, "Big Mac", 5.99, cat.ID, true)

	r := testutils.SetupRouter()
	r.DELETE("/categories/:id", DeleteCategory)

	req := testutils.JSONRequest("DELETE", testutils.IDParam("/categories", cat.ID), nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestDeleteCategory_NotFound(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.DELETE("/categories/:id", DeleteCategory)

	req := testutils.JSONRequest("DELETE", "/categories/999", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
