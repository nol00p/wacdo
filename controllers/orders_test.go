package controllers

import (
	"net/http"
	"testing"
	"wacdo/config"
	"wacdo/models"
	"wacdo/testutils"

	"github.com/stretchr/testify/assert"
)

func seedOrder(userID uint, status string, customerID *uint) models.Order {
	order := models.Order{
		CreatedByID: userID,
		OrderType:   "counter",
		Status:      status,
		CustomerID:  customerID,
		TotalPrice:  10.0,
	}
	config.DB.Create(&order)
	return order
}

func TestCreateOrder_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	user := testutils.SeedUser(db, "admin", "admin@test.com", "P@ssw0rd", role.ID)
	cat := testutils.SeedCategory(db, "Burgers")
	p := testutils.SeedProduct(db, "Big Mac", 5.99, cat.ID, true)

	r := testutils.SetupRouter()
	r.Use(testutils.AuthMiddleware(int(user.ID), "admin"))
	r.POST("/orders", CreateOrder)

	body := map[string]interface{}{
		"order_type": "counter",
		"order_items": []map[string]interface{}{
			{"product_id": p.ID, "quantity": 2},
		},
	}
	req := testutils.JSONRequest("POST", "/orders", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := testutils.ParseResponse(w)
	assert.Equal(t, "pending", resp["status"])
	assert.Equal(t, 5.99*2, resp["total_price"])
}

func TestCreateOrder_WithMenu(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	user := testutils.SeedUser(db, "admin", "admin@test.com", "P@ssw0rd", role.ID)
	m := testutils.SeedMenu(db, "Big Mac Menu", 9.99, true)

	r := testutils.SetupRouter()
	r.Use(testutils.AuthMiddleware(int(user.ID), "admin"))
	r.POST("/orders", CreateOrder)

	body := map[string]interface{}{
		"order_type": "counter",
		"order_items": []map[string]interface{}{
			{"menu_id": m.ID, "quantity": 1},
		},
	}
	req := testutils.JSONRequest("POST", "/orders", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := testutils.ParseResponse(w)
	assert.Equal(t, 9.99, resp["total_price"])
}

func TestCreateOrder_InvalidOrderType(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	user := testutils.SeedUser(db, "admin", "admin@test.com", "P@ssw0rd", role.ID)
	cat := testutils.SeedCategory(db, "Burgers")
	p := testutils.SeedProduct(db, "Big Mac", 5.99, cat.ID, true)

	r := testutils.SetupRouter()
	r.Use(testutils.AuthMiddleware(int(user.ID), "admin"))
	r.POST("/orders", CreateOrder)

	body := map[string]interface{}{
		"order_type": "invalid",
		"order_items": []map[string]interface{}{
			{"product_id": p.ID, "quantity": 1},
		},
	}
	req := testutils.JSONRequest("POST", "/orders", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateOrder_NoItems(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	user := testutils.SeedUser(db, "admin", "admin@test.com", "P@ssw0rd", role.ID)

	r := testutils.SetupRouter()
	r.Use(testutils.AuthMiddleware(int(user.ID), "admin"))
	r.POST("/orders", CreateOrder)

	body := map[string]interface{}{
		"order_type":  "counter",
		"order_items": []map[string]interface{}{},
	}
	req := testutils.JSONRequest("POST", "/orders", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateOrder_BothProductAndMenu(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	user := testutils.SeedUser(db, "admin", "admin@test.com", "P@ssw0rd", role.ID)
	cat := testutils.SeedCategory(db, "Burgers")
	p := testutils.SeedProduct(db, "Big Mac", 5.99, cat.ID, true)
	m := testutils.SeedMenu(db, "Big Mac Menu", 9.99, true)

	r := testutils.SetupRouter()
	r.Use(testutils.AuthMiddleware(int(user.ID), "admin"))
	r.POST("/orders", CreateOrder)

	body := map[string]interface{}{
		"order_type": "counter",
		"order_items": []map[string]interface{}{
			{"product_id": p.ID, "menu_id": m.ID, "quantity": 1},
		},
	}
	req := testutils.JSONRequest("POST", "/orders", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateOrder_UnavailableProduct(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	user := testutils.SeedUser(db, "admin", "admin@test.com", "P@ssw0rd", role.ID)
	cat := testutils.SeedCategory(db, "Burgers")
	p := testutils.SeedProduct(db, "Big Mac", 5.99, cat.ID, false) // not available

	r := testutils.SetupRouter()
	r.Use(testutils.AuthMiddleware(int(user.ID), "admin"))
	r.POST("/orders", CreateOrder)

	body := map[string]interface{}{
		"order_type": "counter",
		"order_items": []map[string]interface{}{
			{"product_id": p.ID, "quantity": 1},
		},
	}
	req := testutils.JSONRequest("POST", "/orders", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetOrders_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	user := testutils.SeedUser(db, "admin", "admin@test.com", "P@ssw0rd", role.ID)
	seedOrder(user.ID, "pending", nil)

	r := testutils.SetupRouter()
	r.GET("/orders", GetOrders)

	req := testutils.JSONRequest("GET", "/orders", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetOrders_FilterByStatus(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	user := testutils.SeedUser(db, "admin", "admin@test.com", "P@ssw0rd", role.ID)
	seedOrder(user.ID, "pending", nil)
	seedOrder(user.ID, "preparing", nil)

	r := testutils.SetupRouter()
	r.GET("/orders", GetOrders)

	req := testutils.JSONRequest("GET", "/orders?status=pending", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetOrder_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	user := testutils.SeedUser(db, "admin", "admin@test.com", "P@ssw0rd", role.ID)
	order := seedOrder(user.ID, "pending", nil)

	r := testutils.SetupRouter()
	r.GET("/orders/:id", GetOrder)

	req := testutils.JSONRequest("GET", testutils.IDParam("/orders", order.ID), nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetOrder_NotFound(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.GET("/orders/:id", GetOrder)

	req := testutils.JSONRequest("GET", "/orders/999", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateOrderStatus_ValidTransition(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	user := testutils.SeedUser(db, "admin", "admin@test.com", "P@ssw0rd", role.ID)
	order := seedOrder(user.ID, "pending", nil)

	r := testutils.SetupRouter()
	r.PATCH("/orders/:id/status", UpdateOrderStatus)

	body := map[string]string{"status": "preparing"}
	req := testutils.JSONRequest("PATCH", testutils.IDParam("/orders", order.ID)+"/status", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := testutils.ParseResponse(w)
	assert.Equal(t, "preparing", resp["status"])
}

func TestUpdateOrderStatus_InvalidTransition(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	user := testutils.SeedUser(db, "admin", "admin@test.com", "P@ssw0rd", role.ID)
	order := seedOrder(user.ID, "pending", nil)

	r := testutils.SetupRouter()
	r.PATCH("/orders/:id/status", UpdateOrderStatus)

	body := map[string]string{"status": "delivered"} // can't go from pending to delivered
	req := testutils.JSONRequest("PATCH", testutils.IDParam("/orders", order.ID)+"/status", body)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateOrderStatus_FullWorkflow(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	user := testutils.SeedUser(db, "admin", "admin@test.com", "P@ssw0rd", role.ID)
	order := seedOrder(user.ID, "pending", nil)

	r := testutils.SetupRouter()
	r.PATCH("/orders/:id/status", UpdateOrderStatus)

	transitions := []string{"preparing", "prepared", "delivered"}
	for _, status := range transitions {
		body := map[string]string{"status": status}
		req := testutils.JSONRequest("PATCH", testutils.IDParam("/orders", order.ID)+"/status", body)
		w := testutils.PerformRequest(r, req)

		assert.Equal(t, http.StatusOK, w.Code, "transition to %s should succeed", status)
	}
}

func TestCancelOrder_FromPending(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	user := testutils.SeedUser(db, "admin", "admin@test.com", "P@ssw0rd", role.ID)
	order := seedOrder(user.ID, "pending", nil)

	r := testutils.SetupRouter()
	r.PATCH("/orders/:id/cancel", CancelOrder)

	req := testutils.JSONRequest("PATCH", testutils.IDParam("/orders", order.ID)+"/cancel", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := testutils.ParseResponse(w)
	assert.Equal(t, "cancelled", resp["status"])
}

func TestCancelOrder_FromPreparing(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	user := testutils.SeedUser(db, "admin", "admin@test.com", "P@ssw0rd", role.ID)
	order := seedOrder(user.ID, "preparing", nil) // not pending

	r := testutils.SetupRouter()
	r.PATCH("/orders/:id/cancel", CancelOrder)

	req := testutils.JSONRequest("PATCH", testutils.IDParam("/orders", order.ID)+"/cancel", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCancelOrder_NotFound(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.PATCH("/orders/:id/cancel", CancelOrder)

	req := testutils.JSONRequest("PATCH", "/orders/999/cancel", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetOrdersByCustomer_Success(t *testing.T) {
	db := testutils.SetupTestDB()
	role := testutils.SeedRole(db, "admin")
	user := testutils.SeedUser(db, "admin", "admin@test.com", "P@ssw0rd", role.ID)
	cust := testutils.SeedCustomer(db, "John", "0612345678", "john@test.com")
	seedOrder(user.ID, "pending", &cust.ID)

	r := testutils.SetupRouter()
	r.GET("/customers/:id/orders", GetOrdersByCustomer)

	req := testutils.JSONRequest("GET", testutils.IDParam("/customers", cust.ID)+"/orders", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetOrdersByCustomer_CustomerNotFound(t *testing.T) {
	testutils.SetupTestDB()

	r := testutils.SetupRouter()
	r.GET("/customers/:id/orders", GetOrdersByCustomer)

	req := testutils.JSONRequest("GET", "/customers/999/orders", nil)
	w := testutils.PerformRequest(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
