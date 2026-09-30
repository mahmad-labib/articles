package validation

import (
	"errors"
	"fmt"
	"net/mail"
	"unicode"
	"unicode/utf8"

	"main.go/internal/domain"
)

var (
	EmailErr    = errors.New("email")
	Empty       = errors.New("field is required")
	NameLength  = errors.New("name length min is 2 and max is 30")
	EmailUnique = errors.New("this email already exists")
	FieldSpaces = errors.New("field cant conatin spaces")
)

func isValidEmail(email string) bool {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	return addr.Address == email
}

type Validator struct{}

func New() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateUser(item *domain.User) error {
	err1 := v.validateName(item.Name)

	err2 := v.ValidateEmail(item.Email)

	err3 := v.validatePassword(item.Password)

	return errors.Join(err1, err2, err3)
}

func (v *Validator) validateName(name string) error {
	nameCount := utf8.RuneCountInString(name)
	if nameCount == 0 {
		return fmt.Errorf("name %v", Empty)
	}
	if nameCount < 2 || nameCount > 30 {
		return NameLength
	}
	return nil
}

func (v *Validator) ValidateEmail(email string) error {
	for _, r := range email {
		if unicode.IsSpace(r) {
			return fmt.Errorf("%v email cant conatin spaces", EmailErr)
		}
	}

	if res := isValidEmail(email); !res {
		return fmt.Errorf("invalid %v", EmailErr)
	}

	emailCount := utf8.RuneCountInString(email)
	if emailCount == 0 {
		return fmt.Errorf("email %v", Empty)
	}
	if emailCount > 100 {
		return fmt.Errorf("%v length max is 100", EmailErr)
	}

	return nil
}

func (v *Validator) validatePassword(pass string) error {
	passCount := utf8.RuneCountInString(pass)
	if passCount == 0 {
		return fmt.Errorf("password %v", Empty)
	}

	for _, r := range pass {
		if unicode.IsSpace(r) {
			return fmt.Errorf("password %v", FieldSpaces)
		}
	}

	return nil
}

func (v *Validator) ValidateCategory(category *domain.Category) error {
	NameLen := utf8.RuneCountInString(category.Name)
	SlugLen := utf8.RuneCountInString(category.Slug)

	if NameLen == 0 || NameLen == 1 {
		return fmt.Errorf("category name %v %v", Empty, " and the min is 2")
	}

	if SlugLen == 0 || SlugLen == 1 {
		return fmt.Errorf("category slug %v %v", Empty, " and the min is 2")
	}

	for _, r := range category.Name {
		if unicode.IsSpace(r) {
			return fmt.Errorf("Category Name %v", FieldSpaces)
		}
	}

	for _, r := range category.Slug {
		if unicode.IsSpace(r) {
			return fmt.Errorf("Category Slug %v", FieldSpaces)
		}
	}

	return nil
}

func (v *Validator) ValidateTag(tag *domain.Tag) error {
	NameLen := utf8.RuneCountInString(tag.Name)

	if NameLen == 0 || NameLen == 1 {
		return fmt.Errorf("tag name %v %v", Empty, " and the min is 2")
	}

	for _, r := range tag.Name {
		if unicode.IsSpace(r) {
			return fmt.Errorf("Tag Name %v", FieldSpaces)
		}
	}

	return nil
}
