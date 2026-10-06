package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/mahmad-labib/articles/internal/domain"
	"github.com/mahmad-labib/articles/internal/infra/db"
)

type ArticleCreation struct {
	Tags       []int  `json:"tags"`
	CategoryID int    `json:"categoryId"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	Status     string `json:"status"`
}

func (h *Handler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var userID int64

	_, claims, _ := jwtauth.FromContext(r.Context())
	if val, ok := claims["user_id"].(float64); ok {
		userID = int64(val)
	}

	var body ArticleCreation

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&body); err != nil {
		handleRequestError(w, err)

		return
	}

	var status domain.ArticleStatus

	if body.Status == "published" {
		status = domain.ArticleStatusPublished
	}

	if body.Status == "scheduled" {
		status = domain.ArticleStatusScheduled
	}

	if body.Status == "draft" || body.Status == "" {
		status = domain.ArticleStatusDraft
	}

	article := domain.Article{
		Title:      body.Title,
		Body:       body.Body,
		Status:     status,
		CategoryID: int64(body.CategoryID),
		AuthorID:   userID,
	}

	tags := make([]domain.Tag, len(body.Tags))
	for i, tagID := range body.Tags {
		tags[i] = domain.Tag{ID: int64(tagID)}
	}

	Article, err := h.asvc.CreateArticle(r.Context(), article, tags)
	if err != nil {
		handleServiceError(w, err)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"Article": Article,
	})
}

func (h *Handler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	var userID int64

	_, claims, _ := jwtauth.FromContext(r.Context())
	if val, ok := claims["user_id"].(float64); ok {
		userID = int64(val)
	}

	articleId := chi.URLParam(r, "id")

	val, err := strconv.ParseInt(articleId, 10, 64)
	if err != nil {
		handleRequestError(w, ErrInvalidId)

		return
	}

	var body ArticleCreation

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&body); err != nil {
		handleRequestError(w, err)

		return
	}

	var status domain.ArticleStatus

	if body.Status == "published" {
		status = domain.ArticleStatusPublished
	}

	if body.Status == "scheduled" {
		status = domain.ArticleStatusScheduled
	}

	if body.Status == "draft" || body.Status == "" {
		status = domain.ArticleStatusDraft
	}

	article := domain.Article{
		ID:         val,
		Title:      body.Title,
		Body:       body.Body,
		Status:     status,
		CategoryID: int64(body.CategoryID),
		AuthorID:   userID,
	}

	tags := make([]domain.Tag, len(body.Tags))
	for i, tagID := range body.Tags {
		tags[i] = domain.Tag{ID: int64(tagID)}
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, "userID", userID)

	Article, err := h.asvc.UpdateArticle(ctx, article, tags)
	if err != nil {
		handleServiceError(w, err)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"Article": Article,
	})
}

func (h *Handler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	var userID int64

	_, claims, _ := jwtauth.FromContext(r.Context())
	if val, ok := claims["user_id"].(float64); ok {
		userID = int64(val)
	}

	id := chi.URLParam(r, "id")

	val, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		handleRequestError(w, ErrInvalidId)

		return
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, "userID", userID)

	err = h.asvc.DeleteArticle(ctx, val)
	if err != nil {
		handleServiceError(w, err)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetArticleByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	val, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		handleRequestError(w, ErrInvalidId)

		return
	}

	Article, err := h.asvc.GetArticleByID(r.Context(), val)
	if err != nil {
		handleServiceError(w, err)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"Article": Article,
	})
}

func (h *Handler) ListPublishedArticles(w http.ResponseWriter, r *http.Request) {
	var filter db.ArticleFilter

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&filter); err != nil {
		handleRequestError(w, err)

		return
	}

	articles, err := h.asvc.ListPublishedArticles(r.Context(), filter)
	if err != nil {
		handleServiceError(w, err)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"articles": articles,
	})
}

func (h *Handler) ListProfileArticles(w http.ResponseWriter, r *http.Request) {
	var userID int64

	_, claims, _ := jwtauth.FromContext(r.Context())
	if val, ok := claims["user_id"].(float64); ok {
		userID = int64(val)
	}

	var body Pagination

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&body); err != nil {
		handleRequestError(w, err)

		return
	}

	articles, err := h.asvc.ListProfileArticles(r.Context(), userID, body.Page, body.PageSize)
	if err != nil {
		handleServiceError(w, err)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"artilces": articles,
	})
}
