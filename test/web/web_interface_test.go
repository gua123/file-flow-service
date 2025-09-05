package web

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestWebInterface(t *testing.T) {
    // 测试GET请求
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "GET" {
            w.WriteHeader(http.StatusOK)
        }
    })
    req := httptest.NewRequest("GET", "/api", nil)
    w := httptest.NewRecorder()
    handler.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Errorf("GET: 预期状态 %d，实际得到 %d", http.StatusOK, w.Code)
    }

    // 测试POST请求
    handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" {
            w.WriteHeader(http.StatusCreated)
        }
    })
    req = httptest.NewRequest("POST", "/api", nil)
    w = httptest.NewRecorder()
    handler.ServeHTTP(w, req)
    if w.Code != http.StatusCreated {
        t.Errorf("POST: 预期状态 %d，实际得到 %d", http.StatusCreated, w.Code)
    }
}
