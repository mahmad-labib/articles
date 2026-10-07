package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mahmad-labib/articles/internal/domain"
	"github.com/mahmad-labib/articles/internal/infra/db"
)

type TagRequest struct {
	Name string `json:"name"`
}

func (h *Handler) CreateTag(w http.ResponseWriter, r *http.Request) {
	var body TagRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&body); err != nil {
		handleRequestError(w, err)

		return
	}

	tag := domain.Tag{
		Name: body.Name,
	}

	id, err := h.tsvc.CreateTag(r.Context(), &tag)
	if err != nil {
		handleServiceError(w, err)

		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"id": id})
}

func (h *Handler) ListTags(w http.ResponseWriter, r *http.Request) {
	var body db.Pagination

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&body); err != nil {
		handleRequestError(w, err)

		return
	}

	tags, err := h.tsvc.ListTags(r.Context(), body)
	if err != nil {
		handleServiceError(w, err)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tags": tags,
	})
}

func (h *Handler) GetTagByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	val, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		handleRequestError(w, ErrInvalidId)

		return
	}

	tag, err := h.tsvc.GetTagByID(r.Context(), val)
	if err != nil {
		handleServiceError(w, err)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tag": tag,
	})
}

func (h *Handler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	val, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		handleRequestError(w, ErrInvalidId)

		return
	}

	var body TagRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&body); err != nil {
		handleRequestError(w, err)

		return
	}

	tag := domain.Tag{
		ID:   val,
		Name: body.Name,
	}

	err = h.tsvc.UpdateTag(r.Context(), &tag)
	if err != nil {
		handleServiceError(w, err)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	val, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		handleRequestError(w, ErrInvalidId)

		return
	}

	err = h.tsvc.DeleteTag(r.Context(), val)
	if err != nil {
		handleServiceError(w, err)

		return
	}

	w.WriteHeader(http.StatusOK)
}
