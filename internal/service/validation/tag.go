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
		return fmt.Errorf("%w tag name min charachters is 2", Validation)
	}

	for _, r := range tag.Name {
		if unicode.IsSpace(r) {
			return fmt.Errorf("%w Tag Name spaces are not allowed", Validation)
		}
	}

	return nil
}
