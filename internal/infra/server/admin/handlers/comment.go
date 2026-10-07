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

func (h *Handler) ListCommentsByArticleID(w http.ResponseWriter, r *http.Request) {
	articleIdStr := chi.URLParam(r, "id")
	articleId, err := strconv.ParseInt(articleIdStr, 10, 64)
	if err != nil {
		handleRequestError(w, ErrInvalidId)
		return
	}

	var filter db.Pagination

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&filter); err != nil {
		handleRequestError(w, err)

		return
	}

	comments, err := h.cosvc.ListCommentsByArticleID(r.Context(), filter, articleId)
	if err != nil {
		handleRequestError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"comments": comments,
	})
}

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	articleIdStr := chi.URLParam(r, "id")
	articleId, err := strconv.ParseInt(articleIdStr, 10, 64)
	if err != nil {
		handleRequestError(w, ErrInvalidId)
		return
	}

	var userID int64

	_, claims, _ := jwtauth.FromContext(r.Context())

	if val, ok := claims["user_id"].(float64); ok {
		userID = int64(val)
	}

	var comment *domain.Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		handleRequestError(w, err)
		return
	}

	comment.ArticleID = articleId
	comment.AuthorID = userID

	err = h.cosvc.CreateComment(r.Context(), comment)
	if err != nil {
		handleRequestError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Comment created successfully",
	})
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentIdStr := chi.URLParam(r, "commentId")
	commentId, err := strconv.ParseInt(commentIdStr, 10, 64)
	if err != nil {
		handleRequestError(w, ErrInvalidId)
		return
	}

	var userID int64

	_, claims, _ := jwtauth.FromContext(r.Context())
	if val, ok := claims["user_id"].(float64); ok {
		userID = int64(val)
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, "userID", userID)

	err = h.cosvc.DeleteComment(ctx, commentId)
	if err != nil {
		handleRequestError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Comment deleted successfully",
	})
}
