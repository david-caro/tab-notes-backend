package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/note/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var gotten_note Note
	err := json.Unmarshal(w.Body.Bytes(), &gotten_note)
	if err != nil {
		return
	}
	expected_note := Note{
		Id:       "1",
		Question: "What is the answer to the ultimate question of life, the universe and everything?",
		Answer:   "42",
	}
	assert.Equal(t, expected_note, gotten_note)
}
