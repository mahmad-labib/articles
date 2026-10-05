package validation

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/mahmad-labib/articles/internal/domain"
)

type TagValidator struct{}

func NewTagValidator() *TagValidator {
	return &TagValidator{}
}

func (v *TagValidator) Validate(tag *domain.Tag) error {
	NameLen := utf8.RuneCountInString(tag.Name)

	if NameLen == 0 || NameLen == 1 {
		return fmt.Errorf("tag name %w %v", Empty, " and the min is 2")
	}

	for _, r := range tag.Name {
		if unicode.IsSpace(r) {
			return fmt.Errorf("Tag Name %w", FieldSpaces)
		}
	}

	return nil
}
