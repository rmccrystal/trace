package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"testing"
)

func TestOnScan(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	m := json.
	httptest.NewRequest("POST", "/api/v1/scan", nil)

	OnScan(c)
}