package controllers

import (
	"net/http"
	"testing"
	"wacdo/config"
	"wacdo/models"
	"wacdo/testutils"

	"github.com/stretchr/testify/assert"
)

func seedOptionDirect(productID uint, name, isUnique string) models.ProductOptions {
	opt := models.ProductOptions{ProductID: productID, Name: name, IsUnique: isUnique, IsRequired: true}
	config.DB.Create(&opt)
	return opt
}

func TestCreateOption_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	cat := testutils.SeedCategory(db, "Burgers")
	p := testutils.SeedProduct(db, "Big Mac", 5.99, cat.ID, true)

	r := testutils.SetupRouter()
	r.POST("/options", CreateOption)

	body := map[string]interface{}{
		"product_id": p.ID,
		"name":       "Size",
		"is_unique":  "single",
	}
	req := testutils.JSONRequest("POST", "/options", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := testutils.ParseResponse(w)
	assert.Equal(t, "Size", resp["name"])
}

func TestCreateOption_ProductNotFound(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.POST("/options", CreateOption)

	body := map[string]interface{}{
		"product_id": 999,
		"name":       "Size",
		"is_unique":  "single",
	}
	req := testutils.JSONRequest("POST", "/options", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateOption_InvalidIsUnique(t *testing.T) {
	db := testutils.SetupTestDB()
	cat := testutils.SeedCategory(db, "Burgers")
	p := testutils.SeedProduct(db, "Big Mac", 5.99, cat.ID, true)

	r := testutils.SetupRouter()
	r.POST("/options", CreateOption)

	body := map[string]interface{}{
		"product_id": p.ID,
		"name":       "Size",
		"is_unique":  "invalid",
	}
	req := testutils.JSONRequest("POST", "/options", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetOptions_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	cat := testutils.SeedCategory(db, "Burgers")
	p := testutils.SeedProduct(db, "Big Mac", 5.99, cat.ID, true)
	seedOptionDirect(p.ID, "Size", "single")

	r := testutils.SetupRouter()
	r.GET("/options", GetOptions)

	req := testutils.JSONRequest("GET", "/options", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetOption_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	cat := testutils.SeedCategory(db, "Burgers")
	p := testutils.SeedProduct(db, "Big Mac", 5.99, cat.ID, true)
	opt := seedOptionDirect(p.ID, "Size", "single")

	r := testutils.SetupRouter()
	r.GET("/options/:id", GetOption)

	req := testutils.JSONRequest("GET", testutils.IDParam("/options", opt.ID), nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetOption_NotFound(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.GET("/options/:id", GetOption)

	req := testutils.JSONRequest("GET", "/options/999", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteOption_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	cat := testutils.SeedCategory(db, "Burgers")
	p := testutils.SeedProduct(db, "Big Mac", 5.99, cat.ID, true)
	opt := seedOptionDirect(p.ID, "Size", "single")

	r := testutils.SetupRouter()
	r.DELETE("/options/:id", DeleteOption)

	req := testutils.JSONRequest("DELETE", testutils.IDParam("/options", opt.ID), nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteOption_NotFound(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.DELETE("/options/:id", DeleteOption)

	req := testutils.JSONRequest("DELETE", "/options/999", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetOptionsByProduct_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	cat := testutils.SeedCategory(db, "Burgers")
	p := testutils.SeedProduct(db, "Big Mac", 5.99, cat.ID, true)
	seedOptionDirect(p.ID, "Size", "single")

	r := testutils.SetupRouter()
	r.GET("/options/product/:product_id", GetOptionsByProduct)

	req := testutils.JSONRequest("GET", testutils.IDParam("/options/product", p.ID), nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetOptionsByProduct_ProductNotFound(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.GET("/options/product/:product_id", GetOptionsByProduct)

	req := testutils.JSONRequest("GET", "/options/product/999", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
