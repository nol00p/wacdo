package controllers

import (
	"net/http"
	"testing"
	"wacdo/testutils"

	"github.com/stretchr/testify/assert"
)

func TestCreateProduct_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	cat := testutils.SeedCategory(db, "Burgers")

	r := testutils.SetupRouter()
	r.POST("/products", CreateProduct)

	body := map[string]interface{}{
		"name":        "Big Mac",
		"price":       5.99,
		"category_id": cat.ID,
	}
	req := testutils.JSONRequest("POST", "/products", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := testutils.ParseResponse(w)
	assert.Equal(t, "Big Mac", resp["name"])
}

func TestCreateProduct_CategoryNotFound(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.POST("/products", CreateProduct)

	body := map[string]interface{}{
		"name":        "Big Mac",
		"price":       5.99,
		"category_id": 999,
	}
	req := testutils.JSONRequest("POST", "/products", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateProduct_Duplicate(t *testing.T) {
	db := testutils.SetupTestDB()
	cat := testutils.SeedCategory(db, "Burgers")
	testutils.SeedProduct(db, "Big Mac", 5.99, cat.ID, true)

	r := testutils.SetupRouter()
	r.POST("/products", CreateProduct)

	body := map[string]interface{}{
		"name":        "Big Mac",
		"price":       6.99,
		"category_id": cat.ID,
	}
	req := testutils.JSONRequest("POST", "/products", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetProducts_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	cat := testutils.SeedCategory(db, "Burgers")
	testutils.SeedProduct(db, "Big Mac", 5.99, cat.ID, true)

	r := testutils.SetupRouter()
	r.GET("/products", GetProducts)

	req := testutils.JSONRequest("GET", "/products", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetProduct_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	cat := testutils.SeedCategory(db, "Burgers")
	p := testutils.SeedProduct(db, "Big Mac", 5.99, cat.ID, true)

	r := testutils.SetupRouter()
	r.GET("/products/:id", GetProduct)

	req := testutils.JSONRequest("GET", testutils.IDParam("/products", p.ID), nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := testutils.ParseResponse(w)
	assert.Equal(t, "Big Mac", resp["name"])
}

func TestGetProduct_NotFound(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.GET("/products/:id", GetProduct)

	req := testutils.JSONRequest("GET", "/products/999", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateProduct_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	cat := testutils.SeedCategory(db, "Burgers")
	p := testutils.SeedProduct(db, "Big Mac", 5.99, cat.ID, true)

	r := testutils.SetupRouter()
	r.PUT("/products/:id", UpdateProduct)

	body := map[string]interface{}{
		"name":        "Big Mac Deluxe",
		"price":       7.99,
		"category_id": cat.ID,
	}
	req := testutils.JSONRequest("PUT", testutils.IDParam("/products", p.ID), body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := testutils.ParseResponse(w)
	assert.Equal(t, "Big Mac Deluxe", resp["name"])
}

func TestUpdateProduct_NameConflict(t *testing.T) {
	db := testutils.SetupTestDB()
	cat := testutils.SeedCategory(db, "Burgers")
	testutils.SeedProduct(db, "Whopper", 6.99, cat.ID, true)
	p := testutils.SeedProduct(db, "Big Mac", 5.99, cat.ID, true)

	r := testutils.SetupRouter()
	r.PUT("/products/:id", UpdateProduct)

	body := map[string]interface{}{
		"name":        "Whopper",
		"price":       5.99,
		"category_id": cat.ID,
	}
	req := testutils.JSONRequest("PUT", testutils.IDParam("/products", p.ID), body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestDeleteProduct_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	cat := testutils.SeedCategory(db, "Burgers")
	p := testutils.SeedProduct(db, "Big Mac", 5.99, cat.ID, true)

	r := testutils.SetupRouter()
	r.DELETE("/products/:id", DeleteProduct)

	req := testutils.JSONRequest("DELETE", testutils.IDParam("/products", p.ID), nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteProduct_NotFound(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.DELETE("/products/:id", DeleteProduct)

	req := testutils.JSONRequest("DELETE", "/products/999", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestToggleProductAvailability_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	cat := testutils.SeedCategory(db, "Burgers")
	p := testutils.SeedProduct(db, "Big Mac", 5.99, cat.ID, true)

	r := testutils.SetupRouter()
	r.PATCH("/products/:id/availability", ToggleProductAvailability)

	req := testutils.JSONRequest("PATCH", testutils.IDParam("/products", p.ID)+"/availability", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := testutils.ParseResponse(w)
	assert.Equal(t, false, resp["is_available"])
}

func TestUpdateProductStock_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	cat := testutils.SeedCategory(db, "Burgers")
	p := testutils.SeedProduct(db, "Big Mac", 5.99, cat.ID, true)

	r := testutils.SetupRouter()
	r.PATCH("/products/:id/stock", UpdateProductStock)

	body := map[string]interface{}{"stock_quantity": 50}
	req := testutils.JSONRequest("PATCH", testutils.IDParam("/products", p.ID)+"/stock", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := testutils.ParseResponse(w)
	assert.Equal(t, float64(50), resp["stock_quantity"])
}

func TestGetProductsByCategory_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	cat := testutils.SeedCategory(db, "Burgers")
	testutils.SeedProduct(db, "Big Mac", 5.99, cat.ID, true)

	r := testutils.SetupRouter()
	r.GET("/products/category/:category_id", GetProductsByCategory)

	req := testutils.JSONRequest("GET", testutils.IDParam("/products/category", cat.ID), nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetProductsByCategory_NotFound(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.GET("/products/category/:category_id", GetProductsByCategory)

	req := testutils.JSONRequest("GET", "/products/category/999", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
