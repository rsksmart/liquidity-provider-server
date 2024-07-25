package utils

import (
	"errors"
	"regexp"
)

var (
	lowerCaseRegex        = regexp.MustCompile(".*[a-z].*")
	upperCaseRegex        = regexp.MustCompile(".*[A-Z].*")
	digitRegex            = regexp.MustCompile(".*[0-9].*")
	specialCharacterRegex = regexp.MustCompile(".*[" + regexp.QuoteMeta(" !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~") + "].*")
)

var PasswordComplexityError = errors.New("password does not meet the complexity requirements")

type ValidationRule func(password string) error

// CheckPasswordComplexity see https://owasp.deteact.com/cheat/cheatsheets/Authentication_Cheat_Sheet.html#password-complexity
func CheckPasswordComplexity(password string, rules ...ValidationRule) error {
	for _, rule := range rules {
		if err := rule(password); err != nil {
			return errors.Join(PasswordComplexityError, err)
		}
	}
	return nil
}

func DefaultPasswordValidationRuleset() []ValidationRule {
	const (
		minLength = 10
		maxLength = 128
	)
	return []ValidationRule{
		PasswordLengthRule(minLength, maxLength),
		PasswordLowerCaseRule(),
		PasswordUpperCaseRule(),
		PasswordDigitRule(),
		PasswordSpecialCharRule(),
	}
}

func PasswordLengthRule(min, max int) ValidationRule {
	return func(password string) error {
		if len(password) < min {
			return errors.New("password is too short")
		}
		if len(password) > max {
			return errors.New("password is too long")
		}
		return nil
	}
}

func PasswordLowerCaseRule() ValidationRule {
	return func(password string) error {
		if !lowerCaseRegex.MatchString(password) {
			return errors.New("password must contain at least one lowercase character")
		}
		return nil
	}
}

func PasswordUpperCaseRule() ValidationRule {
	return func(password string) error {
		if !upperCaseRegex.MatchString(password) {
			return errors.New("password must contain at least one uppercase character")
		}
		return nil
	}
}

func PasswordDigitRule() ValidationRule {
	return func(password string) error {
		if !digitRegex.MatchString(password) {
			return errors.New("password must contain at least one digit")
		}
		return nil
	}
}

func PasswordSpecialCharRule() ValidationRule {
	return func(password string) error {
		if !specialCharacterRegex.MatchString(password) {
			return errors.New("password must contain at least one special character")
		}
		return nil
	}
}
