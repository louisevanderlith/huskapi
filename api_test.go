package huskapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func SetupTest() {

}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestBooks_Router(t *testing.T) {
	CreateContext()
	defer Shutdown()
	router := NewAPI(ctx, 10)
	w := performRequest(router, http.MethodGet, "/books")

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBooks_Get(t *testing.T) {
	CreateContext()
	defer Shutdown()
	router := NewAPI(ctx, 10)
	w := performRequest(router, http.MethodGet, "/books")

	assert.Equal(t, http.StatusOK, w.Code)
	t.Log(w.Body.String())
	var response map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	t.Log(response)
	if err != nil {
		t.Fatal(err)
		return
	}
t.Fail()

}

func TestBooks_View(t *testing.T) {
	CreateContext()
	defer Shutdown()
	router := NewAPI(ctx, 10)
	w := performRequest(router, http.MethodGet, "/books/1580043039`4")

	assert.Equal(t, http.StatusOK, w.Code)
	t.Log(w.Body.String())
	var response map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	t.Log(response)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Fail()

}


func TestBooks_Search(t *testing.T) {
	CreateContext()
	defer Shutdown()
/*
	parm := Book{
		ISBN: "9781593277574",
	}
*/
	router := NewAPI(ctx, 10)
	w := performRequest(router, http.MethodGet, "/books/search/A10/" )

	assert.Equal(t, http.StatusOK, w.Code)
	t.Log(w.Body.String())
	var response map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	t.Log(response)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Fail()

}