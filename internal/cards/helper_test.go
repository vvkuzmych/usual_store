package cards

import (
	"strings"
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid email - standard format",
			email:     "user@example.com",
			wantError: false,
		},
		{
			name:      "valid email - with plus sign",
			email:     "user+tag@example.com",
			wantError: false,
		},
		{
			name:      "valid email - with dot in local part",
			email:     "first.last@example.com",
			wantError: false,
		},
		{
			name:      "valid email - with numbers",
			email:     "user123@example456.com",
			wantError: false,
		},
		{
			name:      "valid email - subdomain",
			email:     "user@mail.example.com",
			wantError: false,
		},
		{
			name:      "valid email - underscore",
			email:     "user_name@example.com",
			wantError: false,
		},
		{
			name:      "valid email - hyphen in domain",
			email:     "user@my-example.com",
			wantError: false,
		},
		{
			name:      "empty email - allowed",
			email:     "",
			wantError: false,
		},
		{
			name:      "invalid email - no @ symbol",
			email:     "userexample.com",
			wantError: true,
			errorMsg:  "invalid email format",
		},
		{
			name:      "invalid email - no domain",
			email:     "user@",
			wantError: true,
			errorMsg:  "invalid email format",
		},
		{
			name:      "invalid email - no local part",
			email:     "@example.com",
			wantError: true,
			errorMsg:  "invalid email format",
		},
		{
			name:      "invalid email - missing TLD",
			email:     "user@example",
			wantError: true,
			errorMsg:  "invalid email format",
		},
		{
			name:      "invalid email - spaces",
			email:     "user name@example.com",
			wantError: true,
			errorMsg:  "invalid email format",
		},
		{
			name:      "invalid email - multiple @ symbols",
			email:     "user@@example.com",
			wantError: true,
			errorMsg:  "invalid email format",
		},
		{
			name:      "invalid email - special characters",
			email:     "user!#$%@example.com",
			wantError: true,
			errorMsg:  "invalid email format",
		},
		{
			name:      "invalid email - missing dot in domain",
			email:     "user@examplecom",
			wantError: true,
			errorMsg:  "invalid email format",
		},
		{
			name:      "edge case - starts with dot (simple regex allows)",
			email:     ".user@example.com",
			wantError: false, // Simple regex doesn't catch this
		},
		{
			name:      "edge case - ends with dot (simple regex allows)",
			email:     "user.@example.com",
			wantError: false, // Simple regex doesn't catch this
		},
		{
			name:      "edge case - consecutive dots (simple regex allows)",
			email:     "user..name@example.com",
			wantError: false, // Simple regex doesn't catch this
		},
		{
			name:      "invalid email - too short TLD",
			email:     "user@example.c",
			wantError: true,
			errorMsg:  "invalid email format",
		},
		{
			name:      "valid email - long TLD",
			email:     "user@example.museum",
			wantError: false,
		},
		{
			name:      "valid email - uppercase",
			email:     "USER@EXAMPLE.COM",
			wantError: false,
		},
		{
			name:      "valid email - mixed case",
			email:     "UsEr@ExAmPlE.CoM",
			wantError: false,
		},
		{
			name:      "invalid email - only @",
			email:     "@",
			wantError: true,
			errorMsg:  "invalid email format",
		},
		{
			name:      "invalid email - only domain",
			email:     "example.com",
			wantError: true,
			errorMsg:  "invalid email format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEmail(tt.email)

			if tt.wantError {
				if err == nil {
					t.Errorf("validateEmail(%q) expected error, got nil", tt.email)
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("validateEmail(%q) error = %v, want error containing %q", tt.email, err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("validateEmail(%q) unexpected error: %v", tt.email, err)
				}
			}
		})
	}
}

// Benchmark for email validation
func BenchmarkValidateEmail(b *testing.B) {
	testEmails := []string{
		"user@example.com",
		"invalid-email",
		"user+tag@example.com",
		"",
		"@example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		email := testEmails[i%len(testEmails)]
		_ = validateEmail(email)
	}
}

// Test edge cases
func TestValidateEmailEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		wantError bool
	}{
		{"very long email", strings.Repeat("a", 100) + "@example.com", false},
		{"very long domain", "user@" + strings.Repeat("a", 100) + ".com", false},
		{"unicode characters", "用户@example.com", true},
		{"emoji", "user😀@example.com", true},
		{"tab character", "user\t@example.com", true},
		{"newline character", "user\n@example.com", true},
		{"null byte", "user\x00@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEmail(tt.email)
			if tt.wantError && err == nil {
				t.Errorf("validateEmail(%q) expected error, got nil", tt.email)
			}
			if !tt.wantError && err != nil {
				t.Errorf("validateEmail(%q) unexpected error: %v", tt.email, err)
			}
		})
	}
}
