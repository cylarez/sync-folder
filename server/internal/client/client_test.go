package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegister(t *testing.T) {
	r, err := http.NewRequest("POST", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	c := Register(w, r)
	check := clients[c.Id]
	if check != c {
		t.Errorf("failed to register client")
	}
}

func TestGet(t *testing.T) {
	c := &Client{10, nil, nil, nil}
	clients[c.Id] = c
	check := Get(c.Id)
	if check != c {
		t.Errorf("failed to get client %d", c.Id)
	}
}

func TestGetClientCount(t *testing.T) {
	expected := GetClientCount() + 1
	c := &Client{2, nil, nil, nil}
	clients[c.Id] = c
	val := GetClientCount()
	if val != expected {
		t.Errorf("Wrong client count get %d but expected %d", val, expected)
	}

}
