package auth

import (
	"net/http"
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

func TestGetBearerToken(t *testing.T) {
	userID := uuid.New()
	token, err := MakeJWT(userID, "secret", time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name    string
		headers http.Header
		want    string
		wantErr bool
	}{
		{
			name:    "Successful test",
			headers: http.Header{"Authorization": []string{token}},
			want:    token,
			wantErr: false,
		},
		{
			name:    "Missing header",
			headers: http.Header{},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Invalid token",
			headers: http.Header{"Bearing": []string{"invalid token"}},
			want:    "",
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := GetBearerToken(tc.headers)
			if (err != nil) != tc.wantErr {
				t.Errorf("GetBearerToken() error: %v, expected error: %v", err, tc.wantErr)
			}
			if actual != tc.want {
				t.Errorf("GetBearerToken() token: %v, expected: %v", actual, tc.want)
			}
		})
	}
}
