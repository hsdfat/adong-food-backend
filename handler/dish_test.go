package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetDishes_WithPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	router := gin.Default()
	router.GET("/dishes", GetDishes)

	// Test case 1: Default pagination
	req, _ := http.NewRequest("GET", "/dishes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "data")
	assert.Contains(t, response, "meta")

	meta := response["meta"].(map[string]interface{})
	assert.Equal(t, float64(1), meta["current_page"])
	assert.Equal(t, float64(10), meta["page_size"])

	// Test case 2: Custom pagination
	req, _ = http.NewRequest("GET", "/dishes?page=2&page_size=5", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	json.Unmarshal(w.Body.Bytes(), &response)
	meta = response["meta"].(map[string]interface{})
	assert.Equal(t, float64(2), meta["current_page"])
	assert.Equal(t, float64(5), meta["page_size"])
}

func TestGetDishes_WithSearch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.GET("/dishes", GetDishes)

	req, _ := http.NewRequest("GET", "/dishes?search=gà", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "data")
	assert.Contains(t, response, "meta")
}
