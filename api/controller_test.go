package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/service"
)

func TestStatus(t *testing.T) {
	router := NewController(service.New_CarPool())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/status", nil)
	router.engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"status":"ok"}`, w.Body.String())
}

func TestAPI(t *testing.T) {
	router := NewController(service.New_CarPool())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/cars", strings.NewReader(`
	[
		{ "id": 1, "maxSeats": 4 },
		{ "id": 2, "maxSeats": 6 }
	]`))
	req.Header = map[string][]string{"Content-Type": {"application/json"}}
	router.engine.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/journey", strings.NewReader(`
	{ "id": 1, "passengers": 4 }
	`))
	req.Header = map[string][]string{"Content-Type": {"application/json"}}
	router.engine.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/locate", strings.NewReader("ID=1"))
	req.Header = map[string][]string{"Content-Type": {"application/x-www-form-urlencoded"}}
	router.engine.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"id":1,"maxSeats":4,"availableSeats":0}`, string(w.Body.Bytes()))

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/dropoff", strings.NewReader("ID=1"))
	req.Header = map[string][]string{"Content-Type": {"application/x-www-form-urlencoded"}}
	router.engine.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
