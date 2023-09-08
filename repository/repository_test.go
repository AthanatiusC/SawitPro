package repository

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

var testCase = CreateUserInput{
	Name:     "test",
	Phone:    "test",
	Password: "test",
}

func TestCreateUser(t *testing.T) {
	e := echo.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockRepositoryInterface(ctrl)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	m.EXPECT().CreateUser(c, gomock.Eq(testCase)).Return(testCase)
}
