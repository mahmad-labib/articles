package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mahmad-labib/articles/internal/domain"
)

type CategoryCreation struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Pagination struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Filter   string `json:"filter"`
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var body CategoryCreation

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&body); err != nil {
		handleRequestError(w, err)

		return
	}

	category := domain.Category{
		Name: body.Name,
		Slug: body.Slug,
	}

	id, err := h.csvc.CreateCategory(r.Context(), &category)
	if err != nil {
		handleServiceError(w, err)

		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"id": id})
}

func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	var body Pagination

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&body); err != nil {
		handleRequestError(w, err)

		return
	}

	list, err := h.csvc.ListCategories(r.Context(), body.Page, body.PageSize)
	if err != nil {
		handleServiceError(w, err)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"list": list,
	})
}

func (h *Handler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	val, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		handleRequestError(w, ErrInvalidId)

		return
	}

	category, err := h.csvc.GetCategoryByID(r.Context(), val)
	if err != nil {
		handleServiceError(w, err)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"category": category,
	})
}

func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	var body CategoryCreation

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&body); err != nil {
		handleRequestError(w, err)

		return
	}

	id := chi.URLParam(r, "id")

	val, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		handleRequestError(w, ErrInvalidId)

		return
	}

	category := domain.Category{
		ID:   val,
		Name: body.Name,
		Slug: body.Slug,
	}

	err = h.csvc.UpdateCategory(r.Context(), &category)
	if err != nil {
		handleServiceError(w, err)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	val, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		handleRequestError(w, ErrInvalidId)

		return
	}

	err = h.csvc.DeleteCategory(r.Context(), val)
	if err != nil {
		handleServiceError(w, err)

		return
	}

	w.WriteHeader(http.StatusOK)
}
