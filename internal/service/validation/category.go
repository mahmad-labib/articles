package validation

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/mahmad-labib/articles/internal/domain"
)

type CategoryValidator struct{}

func NewCategoryValidator() *CategoryValidator {
	return &CategoryValidator{}
}

func (v *CategoryValidator) Validate(category *domain.Category) error {
	NameLen := utf8.RuneCountInString(category.Name)
	SlugLen := utf8.RuneCountInString(category.Slug)

	if NameLen == 0 || NameLen == 1 {
		return fmt.Errorf("%v category name the min char is 2", Validation)
	}

	if SlugLen == 0 || SlugLen == 1 {
		return fmt.Errorf("%v category slug the min char is 2", Validation)
	}

	for _, r := range category.Name {
		if unicode.IsSpace(r) {
			return fmt.Errorf("%v Category Name spaces are not allowed", Validation)
		}
	}

	for _, r := range category.Slug {
		if unicode.IsSpace(r) {
			return fmt.Errorf("%v Category Slug spaces are not allowed", Validation)
		}
	}

	return nil
}
