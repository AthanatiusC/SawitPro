package handler

import (
	"os"
	"testing"

	"github.com/AthanatiusC/SawitPro/repository"
	"github.com/golang/mock/gomock"
)

type TestCaseRequest struct {
	FullName    string `json:"full_name"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
}

var testCases = []TestCaseRequest{
	{"", "", ""},                                      // Missing required fields
	{"Jo", "Passw0rd!", "+628%d%d00000000"},           // Full name length less than 3 characters
	{"John Doe", "Pwd!", "+628%d%d00000000"},          // Password length less than 6 characters
	{"John Doe", "password123!", "+628%d%d00000000"},  // Password without a capital letter
	{"John Doe", "Password123", "+628%d%d00000000"},   // Password without a special character
	{"John Doe", "Password123!", "+628%d%d000000000"}, // Phone number length more than 13 characters
	{"John Doe", "Password123!", "52%d%d00000000"},    // Invalid phone number format
}

func TestValidateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repository.NewMockRepositoryInterface(ctrl)
	opts := NewServerOptions{
		Repository: repo,
		Secret:     os.Getenv("SECRET"),
	}
	h := NewServer(opts)

	for _, tc := range testCases {
		errors := h.ValidateUser(tc)
		if len(errors.Messages) == 0 {
			t.Errorf("Unexpected Condition: \n expected: %s\nactual  : %s", tc, errors.Messages)
		}
	}
}
