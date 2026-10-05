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
		return fmt.Errorf("category name %w %v", Empty, " and the min is 2")
	}

	if SlugLen == 0 || SlugLen == 1 {
		return fmt.Errorf("category slug %w %v", Empty, " and the min is 2")
	}

	for _, r := range category.Name {
		if unicode.IsSpace(r) {
			return fmt.Errorf("Category Name %w", FieldSpaces)
		}
	}

	for _, r := range category.Slug {
		if unicode.IsSpace(r) {
			return fmt.Errorf("Category Slug %w", FieldSpaces)
		}
	}

	return nil
}
