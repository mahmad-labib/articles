package validation

import (
	"errors"
	"fmt"
	"net/mail"
	"unicode"
	"unicode/utf8"

	"github.com/mahmad-labib/articles/internal/domain"
)

var (
	Validation  = errors.New("validation error")
	NotFoundErr = errors.New("Not Found")
)

type UserValidator struct{}

func NewUserValidator() *UserValidator {
	return &UserValidator{}
}

func (v *UserValidator) Validate(item *domain.User) error {
	err1 := v.validateName(item.Name)

	err2 := v.ValidateEmail(item.Email)

	err3 := v.validatePassword(item.Password)

	return errors.Join(err1, err2, err3)
}

func (v *UserValidator) ValidateEmail(email string) error {
	for _, r := range email {
		if unicode.IsSpace(r) {
			return fmt.Errorf("%w email spaces are not allowed", Validation)
		}
	}

	if res := v.isValidEmail(email); !res {
		return fmt.Errorf("%v invalid email", Validation)
	}

	emailCount := utf8.RuneCountInString(email)
	if emailCount == 0 {
		return fmt.Errorf("%w email field is required", Validation)
	}

	if emailCount > 100 {
		return fmt.Errorf("%w email length max is 100", Validation)
	}

	return nil
}

func (v *UserValidator) isValidEmail(email string) bool {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	return addr.Address == email
}

func (v *UserValidator) validateName(name string) error {
	nameCount := utf8.RuneCountInString(name)
	if nameCount == 0 {
		return fmt.Errorf("%w name field is required", Validation)
	}

	if nameCount < 2 || nameCount > 30 {
		return fmt.Errorf("%w name length min is 2 and max is 30", Validation)
	}

	return nil
}

func (v *UserValidator) validatePassword(pass string) error {
	passCount := utf8.RuneCountInString(pass)
	if passCount == 0 {
		return fmt.Errorf("%w password is required", Validation)
	}

	for _, r := range pass {
		if unicode.IsSpace(r) {
			return fmt.Errorf("%w passowrd spaces are not allowed", Validation)
		}
	}

	return nil
}
