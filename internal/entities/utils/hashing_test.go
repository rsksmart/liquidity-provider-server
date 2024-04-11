package utils_test

import (
	"encoding/hex"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/argon2"
	"testing"
)

func TestHashAndSaltArgon2(t *testing.T) {
	type caseType struct {
		message  string
		saltSize int64
	}
	cases := []caseType{
		{
			message:  "test",
			saltSize: 16,
		},
		{
			message:  "other test",
			saltSize: 32,
		},
		{
			message:  "bigger test",
			saltSize: 64,
		},
		{
			message:  "ae49dd4f1b44531367e2f1549262b4870f637d6a3590928037a1b1fbe859d412",
			saltSize: 129,
		},
		{
			message:  "'//12casagASFCADHD+++*",
			saltSize: 128,
		},
	}

	for _, testCase := range cases {
		result, salt, err := utils.HashAndSaltArgon2(testCase.message, testCase.saltSize)
		require.NoError(t, err)
		assert.NotEmpty(t, salt)
		saltBytes, err := hex.DecodeString(salt)
		require.NoError(t, err)
		realHash := argon2.IDKey(
			[]byte(testCase.message),
			saltBytes,
			utils.Argon2Time,
			utils.Argon2Memory,
			utils.Argon2Threads,
			utils.Argon2HashSize,
		)
		assert.Equal(t, realHash, result)
		resultUsingSalt, err := utils.HashArgon2(testCase.message, salt)
		require.NoError(t, err)
		assert.Equal(t, realHash, resultUsingSalt)
	}
}

func TestHashAndSaltArgon2_ErrorHandling(t *testing.T) {
	cases := []struct {
		name     string
		value    string
		saltSize int64
		error    string
	}{
		{
			name:     "Illegal salt size",
			value:    test.AnyString,
			saltSize: 15,
			error:    "the recommended size for a salt in argon2 is at least 16 bytes",
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			result, salt, err := utils.HashAndSaltArgon2(testCase.value, testCase.saltSize)
			assert.Empty(t, result)
			assert.Empty(t, salt)
			require.EqualError(t, err, testCase.error)
		})
	}
}

func TestHashArgon2(t *testing.T) {
	type caseType struct {
		message string
		salt    string
		result  []byte
	}
	cases := []caseType{
		{
			message: "test",
			salt:    "8d9b0dcfebe4a9d87d18153b3adc899f",
			result:  []byte{0x69, 0xc6, 0xa, 0xf, 0xae, 0xf1, 0x33, 0x54, 0x5b, 0x19, 0xae, 0x59, 0x83, 0x1a, 0xc8, 0xa6, 0x10, 0xc5, 0x33, 0xba, 0x3b, 0xe, 0x27, 0x1d, 0xaa, 0x75, 0x7d, 0x26, 0x5, 0x31, 0xd, 0x4e},
		},
		{
			message: "some Password!!+123",
			salt:    "df34e4cc5709b01bbee819608a4dca835d49fd72c3a04adb09e4b162385bcc2f",
			result:  []byte{0x2c, 0x2f, 0x1b, 0x90, 0xca, 0x6d, 0xe8, 0x66, 0xfc, 0x7c, 0xfd, 0x77, 0xa1, 0x74, 0x60, 0xb6, 0x78, 0x7d, 0xfb, 0x74, 0xbb, 0xf, 0x92, 0x6f, 0xcb, 0x90, 0x2b, 0x35, 0xaa, 0x1f, 0x70, 0x68},
		},
		{
			message: "other 123 password !?¿-_",
			salt:    "75a02fe64cd2a89c6d237656e98dbcf2b262b7412136be9a92b9a2426fc2d250e7930fa7a7891cfb9c9338f38de6086705801a1ce64d540ea8d3fcb5ab1e3068",
			result:  []byte{0x33, 0xc5, 0x53, 0x4a, 0xbc, 0xd4, 0x7, 0x2, 0x19, 0x57, 0x6d, 0xee, 0xf, 0xed, 0x7f, 0xbd, 0xe4, 0x11, 0xe1, 0x32, 0xb2, 0xa, 0x98, 0x64, 0xa0, 0x1e, 0xcd, 0xfd, 0xa4, 0x75, 0x94, 0xf3},
		},
	}

	for _, testCase := range cases {
		result, err := utils.HashArgon2(testCase.message, testCase.salt)
		require.NoError(t, err)
		assert.Equal(t, testCase.result, result)
	}
}

func TestHashArgon2_ErrorHandling(t *testing.T) {
	cases := []struct {
		name  string
		value string
		salt  string
		error string
	}{
		{
			name:  "Illegal salt size",
			value: test.AnyString,
			salt:  "801a1ce64d540ea8d3fcb5ab1e3068",
			error: "the recommended size for a salt in argon2 is at least 16 bytes",
		},
		{
			name:  "No hex salt",
			value: test.AnyString,
			salt:  "no hex",
			error: "encoding/hex: invalid byte: U+006E 'n'",
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := utils.HashArgon2(testCase.value, testCase.salt)
			assert.Empty(t, result)
			require.EqualError(t, err, testCase.error)
		})
	}
}
