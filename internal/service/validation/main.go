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
	EmailErr       = errors.New("email")
	Empty          = errors.New("field is required")
	NameLength     = errors.New("name length min is 2 and max is 30")
	EmailUnique    = errors.New("this email already exists")
	PasswordSpaces = errors.New("password cant conatin spaces")
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
			return PasswordSpaces
		}
	}

	return nil
}
