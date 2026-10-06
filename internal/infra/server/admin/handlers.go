package admin

import (
	"net/http"
)

type AuthHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type UserHandler interface {
	UserProfile(w http.ResponseWriter, r *http.Request)
}

type CategoryHandler interface {
	CreateCategory(w http.ResponseWriter, r *http.Request)
	ListCategories(w http.ResponseWriter, r *http.Request)
	GetCategoryByID(w http.ResponseWriter, r *http.Request)
	UpdateCategory(w http.ResponseWriter, r *http.Request)
	DeleteCategory(w http.ResponseWriter, r *http.Request)
}

type TagHandler interface {
	CreateTag(w http.ResponseWriter, r *http.Request)
	ListTags(w http.ResponseWriter, r *http.Request)
	GetTagByID(w http.ResponseWriter, r *http.Request)
	UpdateTag(w http.ResponseWriter, r *http.Request)
	DeleteTag(w http.ResponseWriter, r *http.Request)
}

type ArticleHandler interface {
	CreateArticle(w http.ResponseWriter, r *http.Request)
	GetArticleByID(w http.ResponseWriter, r *http.Request)
	ListPublishedArticles(w http.ResponseWriter, r *http.Request)
	ListProfileArticles(w http.ResponseWriter, r *http.Request)
	UpdateArticle(w http.ResponseWriter, r *http.Request)
	DeleteArticle(w http.ResponseWriter, r *http.Request)
}
