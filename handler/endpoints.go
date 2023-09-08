package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/AthanatiusC/SawitPro/generated"
	"github.com/AthanatiusC/SawitPro/repository"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// (GET /user) Get user endpoint, returns user detail by user-id jwt claims
func (s *Server) GetUser(ctx echo.Context, params generated.GetUserParams) error {
	if params.Authorization == nil {
		return ctx.String(http.StatusForbidden, "unauthorized")
	}

	token := s.ValidateJWT(*params.Authorization)
	if !token.Valid {
		return ctx.String(http.StatusForbidden, "unauthorized")
	}

	idClaims, err := s.GetJWTClaims(token, "id")
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "something went wrong"})
	}

	id, err := strconv.Atoi(idClaims)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "something went wrong"})
	}

	user, err := s.Repository.GetUserById(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "something went wrong"})
	}

	return ctx.JSON(http.StatusOK, generated.User{
		FullName:    user.Name,
		PhoneNumber: user.Phone,
	})
}

// (POST /user) Register user endpoint, register new user with valid phone and password
func (s *Server) Register(ctx echo.Context) error {
	var request generated.RegisterJSONRequestBody
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	errors := s.ValidateUser(request)
	if len(errors.Messages) != 0 {
		return ctx.JSON(http.StatusBadRequest, errors)
	}

	// cost 6 = 64 Rounds(2^6=64) process time<~250ms
	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), 6)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "something went wrong"})
	}

	// avoid id increment because of duplicate violation
	phoneNumber := CleanPhoneNumber(request.PhoneNumber)
	user, err := s.Repository.GetUserByPhoneNumber(ctx.Request().Context(), phoneNumber)
	if err != nil && err != sql.ErrNoRows {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "something went wrong"})
	} else if user.Id != 0 {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "phone number is already registered"})
	}

	result, err := s.Repository.CreateUser(ctx.Request().Context(), repository.CreateUserInput{
		Name:     request.FullName,
		Phone:    phoneNumber,
		Password: string(password),
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "something went wrong"})
	}

	return ctx.JSON(http.StatusOK, generated.RegisterResponse{
		Id: result.Id,
	})
}

// (PUT /user) Update user endpoint, edit user data request with valid authentication and request body
func (s *Server) UpdateUser(ctx echo.Context, params generated.UpdateUserParams) error {
	if params.Authorization == nil {
		return ctx.String(http.StatusForbidden, "unauthorized")
	}

	token := s.ValidateJWT(*params.Authorization)
	if !token.Valid {
		return ctx.String(http.StatusForbidden, "unauthorized")
	}

	idClaims, err := s.GetJWTClaims(token, "id")
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "something went wrong"})
	}

	id, err := strconv.Atoi(idClaims)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "something went wrong"})
	}

	user, err := s.Repository.GetUserById(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "something went wrong"})
	}

	var request generated.UpdateUserJSONRequestBody
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	if request.FullName == "" && request.PhoneNumber == "" {
		return ctx.JSON(http.StatusBadRequest, "request cannot be empty")
	}

	if request.PhoneNumber != "" {
		phoneNumber := CleanPhoneNumber(request.PhoneNumber)
		existingUser, err := s.Repository.GetUserByPhoneNumber(ctx.Request().Context(), phoneNumber)
		if err != nil && err != sql.ErrNoRows {
			return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "something went wrong"})
		} else if existingUser.Id != 0 && phoneNumber != user.Phone {
			return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "phone number is already registered"})
		}
		user.Phone = phoneNumber
	}

	if request.FullName != "" {
		user.Name = request.FullName
	}

	result, err := s.Repository.UpdateUserById(ctx.Request().Context(), repository.UpdateUserInput{
		Id:    user.Id,
		Name:  user.Name,
		Phone: user.Phone,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, generated.User{
		FullName:    result.Name,
		PhoneNumber: result.Phone,
	})
}

// (POST /login) User authentication endpoint, returns valid JWT token to user
func (s *Server) Login(ctx echo.Context) error {
	var request generated.LoginJSONRequestBody
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	errors := s.ValidateUser(request)
	if len(errors.Messages) != 0 {
		return ctx.JSON(http.StatusBadRequest, errors)
	}

	phoneNumber := CleanPhoneNumber(request.PhoneNumber)
	user, err := s.Repository.GetUserByPhoneNumber(ctx.Request().Context(), phoneNumber)
	if err != nil && err != sql.ErrNoRows {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "something went wrong"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "incorrect password or phone number"})
	}

	token, err := s.GenerateJWT(user.Id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "something went wrong"})
	}

	return ctx.JSON(http.StatusOK, generated.LoginResponse{
		Id:    user.Id,
		Token: token,
	})
}
