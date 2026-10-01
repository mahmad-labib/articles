package router

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"main.go/internal/domain"
)

type ArticleCreation struct {
	Tags       []int  `json:"tags"`
	CategoryID int    `json:"categoryId"`
	Title      string `json:"title"`
	Slug       string `json:"slug"`
	Body       string `json:"body"`
	Status     string `json:"status"`
}

func (h *UserHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
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
		Title:  body.Title,
		Slug:   body.Slug,
		Body:   body.Body,
		Status: status,
	}

	tags := make([]domain.Tag, len(body.Tags))
	for i, tagID := range body.Tags {
		tags[i] = domain.Tag{ID: int64(tagID)}
	}

	err := h.svc.CreateArticle(r.Context(), article, userID, int64(body.CategoryID), tags)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
