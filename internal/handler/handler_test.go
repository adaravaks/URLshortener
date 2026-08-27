package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/adaravaks/URLshortener/internal/model"
)

func newTestHandler() (*Handler, *fakeLinkRepository) {
	repo := newFakeLinkRepository()
	return NewHandler(repo), repo
}

func TestCreateLink_Success(t *testing.T) {
	h, _ := newTestHandler()

	body := strings.NewReader(`{"url": "https://example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/links", body)
	w := httptest.NewRecorder()

	h.CreateLink(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var got model.Link
	err := json.NewDecoder(w.Body).Decode(&got)
	require.NoError(t, err)
	assert.Equal(t, "https://example.com", got.OriginalURL)
	assert.NotEmpty(t, got.ShortCode)
	assert.NotZero(t, got.ID)
}

func TestCreateLink_MissingURL(t *testing.T) {
	h, _ := newTestHandler()

	body := strings.NewReader(`{"url": ""}`)
	req := httptest.NewRequest(http.MethodPost, "/links", body)
	w := httptest.NewRecorder()

	h.CreateLink(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateLink_InvalidJSON(t *testing.T) {
	h, _ := newTestHandler()

	body := strings.NewReader(`not json`)
	req := httptest.NewRequest(http.MethodPost, "/links", body)
	w := httptest.NewRecorder()

	h.CreateLink(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRedirect_Success(t *testing.T) {
	h, repo := newTestHandler()
	repo.links["abc1234"] = &model.Link{ID: 1, ShortCode: "abc1234", OriginalURL: "https://example.com"}

	req := httptest.NewRequest(http.MethodGet, "/abc1234", nil)
	req.SetPathValue("code", "abc1234")
	w := httptest.NewRecorder()

	h.Redirect(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "https://example.com", w.Header().Get("Location"))
}

func TestRedirect_NotFound(t *testing.T) {
	h, _ := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/doesnotexist", nil)
	req.SetPathValue("code", "doesnotexist")
	w := httptest.NewRecorder()

	h.Redirect(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestStats_Success(t *testing.T) {
	h, repo := newTestHandler()
	repo.links["abc1234"] = &model.Link{ID: 1, ShortCode: "abc1234", OriginalURL: "https://example.com", ClickCount: 5}

	req := httptest.NewRequest(http.MethodGet, "/links/abc1234/stats", nil)
	req.SetPathValue("code", "abc1234")
	w := httptest.NewRecorder()

	h.Stats(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var got model.Link
	err := json.NewDecoder(w.Body).Decode(&got)
	require.NoError(t, err)
	assert.Equal(t, int64(5), got.ClickCount)
}
