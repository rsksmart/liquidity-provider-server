package utils_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	stringShouldBeValidTemplate   = "string %s should be valid"
	stringShouldBeInvalidTemplate = "string %s should be invalid"
)

var specialChars = []string{
	"!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "-", "_", "+", "=", "{", "}",
	"[", "]", ":", ";", "<", ">", ",", ".", "?", "/", "|", "\\", " ", "\"", "'", "`", "~",
}
var uppercaseChars = []string{
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
	"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
}

var lowercaseChars = []string{
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
	"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
}

var numbers = []string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
}

func TestPasswordLengthRule(t *testing.T) {
	rule := utils.PasswordLengthRule(10, 128)
	var validPasswords = []string{"1234567890", "12345678901", "f4bf9f7fcbedaba0392f108c59d8f4a38b3838efb64877380171b54475c2ade8f4bf9f7fcbedaba0392f108c59d8f4a38b3838efb64877380171b54475c2ade8"}
	var invalidPasswords = []string{"", "123456789", "e7930fa7a7891cfb9c9338f38de6086705801a1ce64d540ea8d3fcb5ab1e3068e7930fa7a7891cfb9c9338f38de6086705801a1ce64d540ea8d3fcb5ab1e30681"}

	for _, password := range validPasswords {
		err := rule(password)
		require.NoError(t, err, "password %s should be valid", password)
	}

	for _, password := range invalidPasswords {
		err := rule(password)
		require.Error(t, err, "password %s should be invalid", password)
	}
}

func TestPasswordSpecialCharRule(t *testing.T) {
	var nonSpecialChars []string
	nonSpecialChars = append(nonSpecialChars, uppercaseChars...)
	nonSpecialChars = append(nonSpecialChars, lowercaseChars...)
	nonSpecialChars = append(nonSpecialChars, numbers...)
	nonSpecialChars = append(nonSpecialChars, "°")
	var allowedStrings = []string{"y!es", ".no", "maybe?", "two strings"}
	var notAllowedStrings = []string{"yes", "no", "maybe", "other°"}

	rule := utils.PasswordSpecialCharRule()
	for _, char := range specialChars {
		err := rule(char)
		require.NoError(t, err, "special char %s should be valid", char)
	}
	for _, char := range nonSpecialChars {
		err := rule(char)
		require.Error(t, err, "non-special char %s should be invalid", char)
	}
	for _, str := range allowedStrings {
		err := rule(str)
		require.NoError(t, err, stringShouldBeValidTemplate, str)
	}
	for _, str := range notAllowedStrings {
		err := rule(str)
		require.Error(t, err, stringShouldBeInvalidTemplate, str)
	}
}

func TestPasswordUpperCaseRule(t *testing.T) {
	var allowedStrings = []string{"YES!", "N!O", "MA+YBE", "TWO STRINGS"}
	var notAllowedStrings = []string{"y!es", "no*", "m/aybe", "other"}

	rule := utils.PasswordUpperCaseRule()
	for _, char := range uppercaseChars {
		err := rule(char)
		require.NoError(t, err, "uppercase char %s should be valid", char)
	}
	for _, char := range lowercaseChars {
		err := rule(char)
		require.Error(t, err, "lowercase char %s should be invalid", char)
	}
	for _, str := range allowedStrings {
		err := rule(str)
		require.NoError(t, err, stringShouldBeValidTemplate, str)
	}
	for _, str := range notAllowedStrings {
		err := rule(str)
		require.Error(t, err, stringShouldBeInvalidTemplate, str)
	}
}

func TestPasswordLowerCaseRule(t *testing.T) {
	var allowedStrings = []string{"yes!", "n!o", "ma+ybe", "two strings"}
	var notAllowedStrings = []string{"YES", "NO", "MA+YBE", "TWO STRINGS"}

	rule := utils.PasswordLowerCaseRule()
	for _, char := range lowercaseChars {
		err := rule(char)
		require.NoError(t, err, "lowercase char %s should be valid", char)
	}
	for _, char := range uppercaseChars {
		err := rule(char)
		require.Error(t, err, "uppercase char %s should be invalid", char)
	}
	for _, str := range allowedStrings {
		err := rule(str)
		require.NoError(t, err, stringShouldBeValidTemplate, str)
	}
	for _, str := range notAllowedStrings {
		err := rule(str)
		require.Error(t, err, stringShouldBeInvalidTemplate, str)
	}
}

func TestPasswordDigitRule(t *testing.T) {
	var allowedStrings = []string{"yes1", "n2o", "ma3ybe", "two strings 4"}
	var notAllowedStrings = []string{"yes", "no", "maybe", "other"}

	rule := utils.PasswordDigitRule()
	for _, char := range numbers {
		err := rule(char)
		require.NoError(t, err, "number %s should be valid", char)
	}
	for _, char := range uppercaseChars {
		err := rule(char)
		require.Error(t, err, "uppercase char %s should be invalid", char)
	}
	for _, str := range allowedStrings {
		err := rule(str)
		require.NoError(t, err, stringShouldBeValidTemplate, str)
	}
	for _, str := range notAllowedStrings {
		err := rule(str)
		require.Error(t, err, stringShouldBeInvalidTemplate, str)
	}
}

func TestCheckPasswordComplexity(t *testing.T) {
	rules := utils.DefaultPasswordValidationRuleset()

	noDigitFakePassword := "NoDigitPassword!"
	noSpecialCharFakePassword := "NoSpecialCharPassword1"
	noUpperCaseFakePassword := "nouppercasepassword1!"
	noLowerCaseFakePassword := "NOLOWERCASEPASSWORD1!"
	tooShortFakePassword := "Short1!"
	// this is not a credential of anything, they are two sha256 hashes concatenated + 1 extra char
	// nolint:gosec
	tooLongFakePassword := "E7930fa7a7891cfb9c9338f38de6086705801a1ce64d540ea8d3fcb5ab1e3068e7930fa7a7891cfb9c9338f38de6086705801a1ce64d540ea8d3fcb5ab1e3068!"
	validFakePassword := "ValidPassword1!"

	err := utils.CheckPasswordComplexity(noDigitFakePassword, rules...)
	require.ErrorContains(t, err, "password must contain at least one digit")
	err = utils.CheckPasswordComplexity(noSpecialCharFakePassword, rules...)
	require.ErrorContains(t, err, "password must contain at least one special character")
	err = utils.CheckPasswordComplexity(noUpperCaseFakePassword, rules...)
	require.ErrorContains(t, err, "password must contain at least one uppercase character")
	err = utils.CheckPasswordComplexity(noLowerCaseFakePassword, rules...)
	require.ErrorContains(t, err, "password must contain at least one lowercase character")
	err = utils.CheckPasswordComplexity(tooShortFakePassword, rules...)
	require.ErrorContains(t, err, "password is too short")
	err = utils.CheckPasswordComplexity(tooLongFakePassword, rules...)
	require.ErrorContains(t, err, "password is too long")
	err = utils.CheckPasswordComplexity(validFakePassword, rules...)
	require.NoError(t, err)
	require.Len(t, rules, 5)
}
