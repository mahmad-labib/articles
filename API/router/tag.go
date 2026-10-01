package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"main.go/internal/domain"
)

type TagRequest struct {
	Name string `json:"name"`
}

func (h *UserHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
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

	id, err := h.svc.CreateTag(r.Context(), &tag)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"id": id})
}

func (h *UserHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	var body Pagination

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&body); err != nil {
		handleRequestError(w, err)
		return
	}

	tags, err := h.svc.ListTags(r.Context(), body.Page, body.PageSize)
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

func (h *UserHandler) GetTagByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	val, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		handleRequestError(w, ErrInvalidId)
		return
	}

	tag, err := h.svc.GetTagByID(r.Context(), val)
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

func (h *UserHandler) UpdateTag(w http.ResponseWriter, r *http.Request) {
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

	err = h.svc.UpdateTag(r.Context(), &tag)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	val, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		handleRequestError(w, ErrInvalidId)
		return
	}

	err = h.svc.DeleteTag(r.Context(), val)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
