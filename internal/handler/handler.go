package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/adaravaks/URLshortener/internal/model"
	"github.com/adaravaks/URLshortener/internal/repository"
)

type Handler struct {
	repo repository.LinkRepository
}

func NewHandler(repo repository.LinkRepository) *Handler {
	return &Handler{repo: repo}
}

type createLinkRequest struct {
	URL string `json:"url"`
}

func (h *Handler) CreateLink(w http.ResponseWriter, r *http.Request) {
	var req createLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.URL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	code, err := generateShortCode()
	if err != nil {
		http.Error(w, "failed to generate code", http.StatusInternalServerError)
		return
	}

	link := &model.Link{
		ShortCode:   code,
		OriginalURL: req.URL,
	}

	if err := h.repo.Create(r.Context(), link); err != nil {
		http.Error(w, "failed to create link", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(link); err != nil {
		http.Error(w, "failed to encode to json", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	link, err := h.repo.GetByCode(r.Context(), code)
	if errors.Is(err, repository.ErrNotFound) {
		http.Error(w, "link not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to look up link", http.StatusInternalServerError)
		return
	}

	go func() {
		ctx := context.Background()
		err := h.repo.IncrementClickCount(ctx, code)
		if err != nil {
			log.Println("failed to increment click count:", err)
		}
	}()

	http.Redirect(w, r, link.OriginalURL, http.StatusFound)
}

func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	link, err := h.repo.GetByCode(r.Context(), code)
	if errors.Is(err, repository.ErrNotFound) {
		http.Error(w, "link not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to look up link", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(link); err != nil {
		http.Error(w, "failed to encode to json", http.StatusInternalServerError)
		return
	}
}
