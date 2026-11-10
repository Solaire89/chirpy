package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	type testCase struct {
		name        string
		password    string
		hash        string
		expectMatch bool
		expectError bool
	}
	someHash, err := HashPassword("somepassword")
	if err != nil {
		t.Fatal(err)
	}
	badHash := "badhash"
	cases := []testCase{
		{
			name:        "correct password matches",
			password:    "somepassword",
			hash:        someHash,
			expectMatch: true,
			expectError: false,
		},
		{
			name:        "incorrect password match",
			password:    "differentpassword",
			hash:        someHash,
			expectMatch: false,
			expectError: false,
		},
		{
			name:        "bad hash",
			password:    "somepassword",
			hash:        badHash,
			expectMatch: false,
			expectError: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := CheckPasswordHash(tc.password, tc.hash)
			matchWrong := actual != tc.expectMatch
			errorWrong := (err != nil) != tc.expectError
			if matchWrong || errorWrong {
				t.Errorf("CheckPasswordHash(%q, %q) = (%v %v); expected match: %v, expected error: %v",
					tc.password, tc.hash, actual, err, tc.expectMatch, tc.expectError)
			}
		})
	}
	hashOne, err := HashPassword("apassword")
	if err != nil {
		t.Fatal(err)
	}
	hashTwo, err := HashPassword("apassword")
	if err != nil {
		t.Fatal(err)
	}
	if hashOne == hashTwo {
		t.Errorf("hashes match: %s = %s", hashOne, hashTwo)
	}
}
