package quote_test

import (
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
)

func TestValidateQuoteHash(t *testing.T) {
	tests := []struct {
		name    string
		hash    string
		wantErr bool
	}{
		{
			name:    "Valid 32-byte hash",
			hash:    "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			wantErr: false,
		},
		{
			name:    "Invalid length - too short",
			hash:    "1234567890abcdef",
			wantErr: true,
		},
		{
			name:    "Invalid length - too long",
			hash:    "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef00",
			wantErr: true,
		},
		{
			name:    "Invalid characters",
			hash:    "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdeg",
			wantErr: true,
		},
		{
			name:    "Empty string",
			hash:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := quote.ValidateQuoteHash(tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateQuoteHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
