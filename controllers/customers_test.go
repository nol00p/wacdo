package controllers

import (
	"net/http"
	"testing"
	"wacdo/testutils"

	"github.com/stretchr/testify/assert"
)

func TestCreateCustomer_Success(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.POST("/customers", CreateCustomer)

	body := map[string]string{"name": "John Doe", "phone": "0612345678", "email": "john@test.com"}
	req := testutils.JSONRequest("POST", "/customers", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := testutils.ParseResponse(w)
	assert.Equal(t, "John Doe", resp["name"])
}

func TestCreateCustomer_DuplicatePhone(t *testing.T) {
	db := testutils.SetupTestDB()
	testutils.SeedCustomer(db, "Jane", "0612345678", "jane@test.com")

	r := testutils.SetupRouter()
	r.POST("/customers", CreateCustomer)

	body := map[string]string{"name": "John Doe", "phone": "0612345678"}
	req := testutils.JSONRequest("POST", "/customers", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestCreateCustomer_InvalidData(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.POST("/customers", CreateCustomer)

	// Missing required "name" field
	body := map[string]string{"phone": "0612345678"}
	req := testutils.JSONRequest("POST", "/customers", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetCustomers_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	testutils.SeedCustomer(db, "John", "0600000001", "john@test.com")
	testutils.SeedCustomer(db, "Jane", "0600000002", "jane@test.com")

	r := testutils.SetupRouter()
	r.GET("/customers", GetCustomers)

	req := testutils.JSONRequest("GET", "/customers", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetCustomer_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	c := testutils.SeedCustomer(db, "John", "0612345678", "john@test.com")

	r := testutils.SetupRouter()
	r.GET("/customers/:id", GetCustomer)

	req := testutils.JSONRequest("GET", testutils.IDParam("/customers", c.ID), nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := testutils.ParseResponse(w)
	assert.Equal(t, "John", resp["name"])
}

func TestGetCustomer_NotFound(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.GET("/customers/:id", GetCustomer)

	req := testutils.JSONRequest("GET", "/customers/999", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateCustomer_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	c := testutils.SeedCustomer(db, "John", "0612345678", "john@test.com")

	r := testutils.SetupRouter()
	r.PUT("/customers/:id", UpdateCustomer)

	body := map[string]string{"name": "John Updated", "phone": "0699999999"}
	req := testutils.JSONRequest("PUT", testutils.IDParam("/customers", c.ID), body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateCustomer_PhoneConflict(t *testing.T) {
	db := testutils.SetupTestDB()
	testutils.SeedCustomer(db, "Jane", "0699999999", "jane@test.com")
	c := testutils.SeedCustomer(db, "John", "0612345678", "john@test.com")

	r := testutils.SetupRouter()
	r.PUT("/customers/:id", UpdateCustomer)

	body := map[string]string{"name": "John", "phone": "0699999999"}
	req := testutils.JSONRequest("PUT", testutils.IDParam("/customers", c.ID), body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestDeleteCustomer_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	c := testutils.SeedCustomer(db, "John", "0612345678", "john@test.com")

	r := testutils.SetupRouter()
	r.DELETE("/customers/:id", DeleteCustomer)

	req := testutils.JSONRequest("DELETE", testutils.IDParam("/customers", c.ID), nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteCustomer_NotFound(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.DELETE("/customers/:id", DeleteCustomer)

	req := testutils.JSONRequest("DELETE", "/customers/999", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
