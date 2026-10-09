package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mahmad-labib/articles/internal/domain"
)

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	tags := make(chan *[]domain.Tag, 1)
	categories := make(chan *[]domain.Category, 1)
	articles := make(chan *[]domain.Article, 1)

	go func() {
		tags <- h.SearchTags(r.Context(), name)
	}()
	go func() {
		categories <- h.SearchCategories(r.Context(), name)
	}()
	go func() {
		articles <- h.SearchArticles(r.Context(), name)
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tags":       <-tags,
		"categories": <-categories,
		"articles":   <-articles,
	})
}

func (h *Handler) SearchTags(ctx context.Context, name string) *[]domain.Tag {
	tags, _ := h.tsvc.GetTagsByName(ctx, name)
	return tags
}

func (h *Handler) SearchCategories(ctx context.Context, slug string) *[]domain.Category {
	categories, _ := h.csvc.GetCategoryListBySlug(ctx, slug)
	return categories
}

func (h *Handler) SearchArticles(ctx context.Context, title string) *[]domain.Article {
	articles, _ := h.asvc.GetArticleListByTitle(ctx, title)
	return articles
}
