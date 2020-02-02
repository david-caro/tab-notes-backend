package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNoteRoute(t *testing.T) {
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
		Question: "What is the answer to the ultimate question of life, the universe and everything?",
		Answer:   "42",
	}
	assert.Equal(t, expected_note, gotten_note)
}

func TestPostNoteRoute(t *testing.T) {
	router := SetupRouter()
	expected_question := "What is the answer to the ultimate question of life, the universe and everything?"
	expected_answer := "42"
	expected_note := Note{
		Question: expected_question,
		Answer:   expected_answer,
	}
	request_body := fmt.Sprintf(
		"{'question': '%v', 'answer': '%v'}",
		expected_note.Question,
		expected_note.Answer,
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"POST",
		"/note/1",
		strings.NewReader(request_body),
	)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var gotten_note Note
	err := json.Unmarshal(w.Body.Bytes(), &gotten_note)
	if err != nil {
		return
	}

	assert.Equal(t, expected_note, gotten_note)
}
