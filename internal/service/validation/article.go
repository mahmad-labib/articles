package validation

import (
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/mahmad-labib/articles/internal/domain"
)

type ArticleValidator struct{}

func NewArticleValidator() *ArticleValidator {
	return &ArticleValidator{}
}

func (v *ArticleValidator) Validate(article *domain.Article) error {
	err1 := v.validateTitle(article.Title)
	err2 := v.validateBody(article.Body)

	return errors.Join(err1, err2)
}

func (v *ArticleValidator) validateTitle(name string) error {
	nameCount := utf8.RuneCountInString(name)
	if nameCount == 0 {
		return fmt.Errorf("%w title field is required", Validation)
	}

	if nameCount < 2 || nameCount > 30 {
		return fmt.Errorf("%w title length min is 2 and max is 30", Validation)
	}

	return nil
}

func (v *ArticleValidator) validateBody(name string) error {
	nameCount := utf8.RuneCountInString(name)
	if nameCount == 0 {
		return fmt.Errorf("%w Body is required", Validation)
	}

	if nameCount < 100 || nameCount > 4000 {
		return fmt.Errorf("%w Body length min is 100 and max is 4000", Validation)
	}

	return nil
}
