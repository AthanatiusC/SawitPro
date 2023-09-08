package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/AthanatiusC/SawitPro/generated"
	"github.com/AthanatiusC/SawitPro/repository"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// Since Bcrypt value is always salted and different result, we create custom matcher for it
type RegisterValidator struct {
	FullName    string
	PhoneNumber string
	Password    string
}

func (r RegisterValidator) Matches(values interface{}) bool {
	user, ok := values.(repository.CreateUserInput)
	if ok {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.Password)); err != nil {
			return false
		}
		if r.FullName != user.Name || r.PhoneNumber != user.Phone {
			return false
		}
	}
	return true
}

func (r RegisterValidator) String() string {
	return r.Password
}

// Custom Matcher for GoMock validation on register
func ValidateCreateUser(user repository.CreateUserInput) gomock.Matcher {
	return RegisterValidator{user.Name, user.Phone, user.Password}
}

/*
TestGetUser Criteria:
- Valid JWT
- Repository returns match value
- Assert no error on call
*/
func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	repo := repository.NewMockRepositoryInterface(ctrl)
	opts := NewServerOptions{Repository: repo}
	h := NewServer(opts)

	var userId int = 1 // default user
	repo.EXPECT().GetUserById(gomock.Any(), userId).Return(repository.User{Id: 0, Name: "user"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	token, err := h.GenerateJWT(userId)
	if err != nil {
		t.Error(err)
	}
	token = fmt.Sprintf("Bearer %s", token)

	rec := httptest.NewRecorder()
	if assert.NoError(t, h.GetUser(e.NewContext(req, rec), generated.GetUserParams{Authorization: &token})) {
		assert.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	}

}

/*
TestRegister Criteria:
- Valid User Request
- Repository returns match value
- Repository returns valid bcrypt
- Assert no error on call
*/
func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	repo := repository.NewMockRepositoryInterface(ctrl)
	opts := NewServerOptions{
		Repository: repo,
		Secret:     os.Getenv("SECRET"),
	}
	h := NewServer(opts)

	request := TestCaseRequest{
		FullName:    "John Doe",
		Password:    "Passw0rd!",
		PhoneNumber: "+628%d%d00000000",
	}

	request.PhoneNumber = fmt.Sprintf(request.PhoneNumber, rand.Intn(9), rand.Intn(9))
	validRPN := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(request.PhoneNumber, "")
	user := repository.User{
		Name:      request.FullName,
		Phone:     validRPN,
		Password:  request.Password,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	repo.EXPECT().GetUserByPhoneNumber(gomock.Any(), user.Phone).Return(repository.User{}, nil)
	repo.EXPECT().CreateUser(gomock.Any(), ValidateCreateUser(repository.CreateUserInput{ //  Custom validator for Bcrypt
		Name:     user.Name,
		Phone:    user.Phone,
		Password: user.Password,
	})).Return(user, nil)

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	if assert.NoError(t, h.Register(e.NewContext(req, rec))) {
		assert.Equal(t, http.StatusOK, rec.Code, rec.Body.String(), string(jsonRequest))
	}
}

/*
TestRegister Criteria:
- Valid User Request
- Assert no double phone number
- Repository returns updated value
- Assert no error on call
*/
func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	repo := repository.NewMockRepositoryInterface(ctrl)
	opts := NewServerOptions{
		Repository: repo,
	}
	h := NewServer(opts)

	request := TestCaseRequest{
		FullName:    "user",
		PhoneNumber: fmt.Sprintf("+628%d%d00000000", rand.Intn(9), rand.Intn(9)),
	}

	response := repository.User{
		Id:    1,
		Name:  request.FullName,
		Phone: request.PhoneNumber,
	}

	repo.EXPECT().GetUserByPhoneNumber(gomock.Any(), CleanPhoneNumber(request.PhoneNumber)).Return(repository.User{}, nil)
	repo.EXPECT().GetUserById(gomock.Any(), response.Id).Return(response, nil)
	repo.EXPECT().UpdateUserById(gomock.Any(), repository.UpdateUserInput{
		Id:    1,
		Name:  request.FullName,
		Phone: CleanPhoneNumber(request.PhoneNumber),
	}).Return(response, nil)

	token, err := h.GenerateJWT(response.Id)
	if err != nil {
		t.Error(err)
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		t.Error(jsonRequest)
	}

	req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	token = fmt.Sprintf("Bearer %s", token)
	rec := httptest.NewRecorder()
	if assert.NoError(t, h.UpdateUser(e.NewContext(req, rec), generated.UpdateUserParams{Authorization: &token})) {
		assert.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	}
}

/*
TestRegister Criteria:
- Valid User Request
- Assert phone number and password combination match
- Assert bcrypt password valid
- Assert no error on call
*/
func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	repo := repository.NewMockRepositoryInterface(ctrl)
	opts := NewServerOptions{
		Repository: repo,
	}
	h := NewServer(opts)

	request := generated.LoginJSONRequestBody{
		PhoneNumber: "+6280000000000",
		Password:    "Userpassw0rd!",
	}

	repo.EXPECT().GetUserByPhoneNumber(gomock.Any(), CleanPhoneNumber(request.PhoneNumber)).Return(repository.User{Password: "$2a$06$bt380.sYY0HEAa1tz2eyfOOQDHarjgiABmv.ZJTXzKdXMU.hQFAyi"}, nil)

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		t.Error(t)
	}

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	if assert.NoError(t, h.Login(e.NewContext(req, rec))) {
		assert.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	}
}
