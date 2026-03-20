package controllers

import (
	"net/http"
	"testing"
	"wacdo/config"
	"wacdo/models"
	"wacdo/testutils"

	"github.com/stretchr/testify/assert"
)

func seedOptionValue(optionID uint, value string, price float64) models.OptionValues {
	ov := models.OptionValues{OptionID: optionID, Value: value, OptionPrice: price}
	config.DB.Create(&ov)
	return ov
}

func TestCreateOptionValue_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	cat := testutils.SeedCategory(db, "Burgers")
	p := testutils.SeedProduct(db, "Big Mac", 5.99, cat.ID, true)
	opt := seedOptionDirect(p.ID, "Size", "single")

	r := testutils.SetupRouter()
	r.POST("/options/:id/values", CreateOptionValue)

	body := []map[string]interface{}{
		{"value": "Small", "option_price": 0.0},
		{"value": "Large", "option_price": 1.5},
	}
	req := testutils.JSONRequest("POST", testutils.IDParam("/options", opt.ID)+"/values", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateOptionValue_OptionNotFound(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.POST("/options/:id/values", CreateOptionValue)

	body := []map[string]interface{}{{"value": "Small", "option_price": 0.0}}
	req := testutils.JSONRequest("POST", "/options/999/values", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateOptionValue_DuplicateValue(t *testing.T) {
	db := testutils.SetupTestDB()
	cat := testutils.SeedCategory(db, "Burgers")
	p := testutils.SeedProduct(db, "Big Mac", 5.99, cat.ID, true)
	opt := seedOptionDirect(p.ID, "Size", "single")
	seedOptionValue(opt.ID, "Small", 0)

	r := testutils.SetupRouter()
	r.POST("/options/:id/values", CreateOptionValue)

	body := []map[string]interface{}{{"value": "Small", "option_price": 0.0}}
	req := testutils.JSONRequest("POST", testutils.IDParam("/options", opt.ID)+"/values", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetOptionValue_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	cat := testutils.SeedCategory(db, "Burgers")
	p := testutils.SeedProduct(db, "Big Mac", 5.99, cat.ID, true)
	opt := seedOptionDirect(p.ID, "Size", "single")
	ov := seedOptionValue(opt.ID, "Small", 0)

	r := testutils.SetupRouter()
	r.GET("/options/values/:id", GetOptionValue)

	req := testutils.JSONRequest("GET", testutils.IDParam("/options/values", ov.ID), nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := testutils.ParseResponse(w)
	assert.Equal(t, "Small", resp["value"])
}

func TestGetOptionValue_NotFound(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.GET("/options/values/:id", GetOptionValue)

	req := testutils.JSONRequest("GET", "/options/values/999", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteOptionValue_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	cat := testutils.SeedCategory(db, "Burgers")
	p := testutils.SeedProduct(db, "Big Mac", 5.99, cat.ID, true)
	opt := seedOptionDirect(p.ID, "Size", "single")
	ov := seedOptionValue(opt.ID, "Small", 0)

	r := testutils.SetupRouter()
	r.DELETE("/options/values/:id", DeleteOptionValue)

	req := testutils.JSONRequest("DELETE", testutils.IDParam("/options/values", ov.ID), nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetValuesByOption_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	cat := testutils.SeedCategory(db, "Burgers")
	p := testutils.SeedProduct(db, "Big Mac", 5.99, cat.ID, true)
	opt := seedOptionDirect(p.ID, "Size", "single")
	seedOptionValue(opt.ID, "Small", 0)
	seedOptionValue(opt.ID, "Large", 1.5)

	r := testutils.SetupRouter()
	r.GET("/options/:id/values", GetValuesByOption)

	req := testutils.JSONRequest("GET", testutils.IDParam("/options", opt.ID)+"/values", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetValuesByOption_OptionNotFound(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.GET("/options/:id/values", GetValuesByOption)

	req := testutils.JSONRequest("GET", "/options/999/values", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
