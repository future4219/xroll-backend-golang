package validation

import (
	"errors"
	"regexp"
	"unicode/utf8"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entconst"
)

var (
	ErrInvalidPhoneNumberLength = errors.New("invalid phone number length")
	ErrInvalidPhoneNumber       = errors.New("invalid phone number")
	ErrInvalidCharacterLimit    = errors.New("invalid character limit")
	ErrInvalidCharacterLength   = errors.New("invalid character length")
)

func ValidatePhoneNumber(phone string) error {
	pLength := len(phone)
	if pLength != 10 && pLength != 11 {
		return ErrInvalidPhoneNumberLength
	}

	const pattern = `^0\d{9,10}$`
	var phoneRegexp *regexp.Regexp = regexp.MustCompile(pattern)
	if !phoneRegexp.MatchString(phone) {
		return ErrInvalidPhoneNumber
	}

	return nil
}

type LimitOptions struct {
	limit int
}

type LimitOption func(*LimitOptions)

func Limit(arg1 int) LimitOption {
	return func(opts *LimitOptions) {
		opts.limit = arg1
	}
}

func ValidateCharacterLimit(s string, options ...LimitOption) error {
	// defaultの文字数は255
	opts := &LimitOptions{
		limit: 255,
	}

	for _, option := range options {
		option(opts)
	}

	if utf8.RuneCountInString(s) > opts.limit {
		return entconst.NewValidationError(ErrInvalidCharacterLimit.Error())
	}

	return nil
}

func ValidateCharacterLength(s string, length int) error {
	if s != "" && utf8.RuneCountInString(s) != length {
		return ErrInvalidCharacterLength
	}

	return nil
}

func ValidateKataKana(s string) error {
	const pattern = `^[ァ-ヶー]+$`
	var hiraganaRegexp *regexp.Regexp = regexp.MustCompile(pattern)
	if !hiraganaRegexp.MatchString(s) {
		return entconst.NewValidationError("invalid hiragana")
	}

	return nil
}
