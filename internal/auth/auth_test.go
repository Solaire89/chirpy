package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	successfulToken, err := MakeJWT(userID, "secretToken", time.Hour)
	if err != nil {
		t.Errorf("error creating token: %v", err)
	}
	testCase := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		expectError bool
	}{
		{
			name:        "Valid token",
			tokenString: successfulToken,
			tokenSecret: "secretToken",
			wantUserID:  userID,
			expectError: false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secretToken",
			wantUserID:  uuid.Nil,
			expectError: true,
		},
		{
			name:        "Invalid token secret",
			tokenString: successfulToken,
			tokenSecret: "invalid",
			wantUserID:  uuid.Nil,
			expectError: true,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ValidateJWT(tc.tokenString, tc.tokenSecret)
			if (err != nil) != tc.expectError {
				t.Errorf("ValidateJWT() error = %v; expectError: %v", err, tc.expectError)
				return
			}
			if actual != tc.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", actual, tc.wantUserID)
			}
		})
	}
}
