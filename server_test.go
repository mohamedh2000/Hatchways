package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApiPost(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "api/ping?tags=tech", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(ApiPost)
	handler.ServeHTTP(res, req)

	fmt.Println(res.Body.String())

	if http.StatusOK != res.Code {
		t.Error("call failed")
	}
}

func TestApiPostsErrorTag(t *testing.T) {
	noTagMsg := `{"error":"Tags parameter is required"}`

	illTag, _ := http.NewRequest(http.MethodGet, "/api/ping", nil)
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(ApiPost)
	handler.ServeHTTP(res, illTag)

	if res.Body.String() != noTagMsg {
		t.Error("Should not have passed")
	}
}

func TestApiPostsErrorDirection(t *testing.T) {
	illegalDirection := `{"error":"direction parameter is invalid"}`

	illTag, _ := http.NewRequest(http.MethodGet, "/api/ping?tags=tech&direction=temp", nil)
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(ApiPost)
	handler.ServeHTTP(res, illTag)

	if res.Body.String() != illegalDirection {
		t.Error("Should not have passed")
	}

}

func TestApiPostsErrorSort(t *testing.T) {
	illegalSort := `{"error":"sortBy parameter is invalid"}`

	illSort, _ := http.NewRequest(http.MethodGet, "/api/ping?tags=tech&sortBy=temp", nil)
	resII := httptest.NewRecorder()
	handlerIV := http.HandlerFunc(ApiPost)
	handlerIV.ServeHTTP(resII, illSort)
	if resII.Body.String() != illegalSort {
		t.Error("Should not have passed")
	}
}

func TestApiPing(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/api/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(ApiPing)
	handler.ServeHTTP(res, req)

	if http.StatusOK != res.Code {
		t.Error("call failed")
	}
	if res.Body.String() != `{"success":true}` {
		t.Error(`Expected {"success":true}, got`, res.Body.String())
	}
}
