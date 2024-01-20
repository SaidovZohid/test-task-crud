package v1

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/SaidovZohid/test-task-crud/api/models"
	"github.com/SaidovZohid/test-task-crud/pkg/email"
	"github.com/SaidovZohid/test-task-crud/pkg/utils"
	"github.com/SaidovZohid/test-task-crud/storage/repo"
	"github.com/gin-gonic/gin"
)

// @Router /auth/signup [post]
// @Summary Sign up api for new users
// @Description Sign up to the application with email and password and get the verification code and verify email address.
// @Tags auth
// @Accept json
// @Produce json
// @Param data body models.AuthRequest true "Data"
// @Success 200 {object} models.Message
// @Failure 500 {object} models.ResponseError
// @Failure 400 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) SignUp(ctx *gin.Context) {
	var req models.AuthRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	_, err = h.storage.User().GetByEmail(ctx.Request.Context(), req.Email)
	if !errors.Is(err, sql.ErrNoRows) {
		h.log.Error(err)
		ctx.JSON(http.StatusBadRequest, errResponse(ErrEmailExists))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		h.log.Error(err)
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	req.Password = hashedPassword
	userData, err := json.Marshal(req)
	if err != nil {
		h.log.Error(err)
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	if err != nil {
		h.log.Error(err)
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	code, err := utils.GenerateRandomCode(6)
	if err != nil {
		h.log.Error(err)
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	if err := h.inMemory.Set("sign_up_"+req.Email, string(userData), h.cfg.SignUpAuthicationTime); err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	if err := h.inMemory.Set("code_"+req.Email, code, h.cfg.SignUpAuthicationTime); err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	go func() {
		if err := email.SendEmail(h.cfg, &email.SendEmailRequest{
			To:   []string{req.Email},
			Type: email.VerificationEmail,
			Body: map[string]string{
				"code": code,
				"name": "person who is testing my application :)",
			},
			Subject: "Verification code by zohiddev.me",
		}); err != nil {
			h.log.Error(err)
		}
	}()

	ctx.JSON(http.StatusOK, models.Message{
		Msg: fmt.Sprintf("Check your email for a validation code sent to \"%s\". Use the code to verify your email address.", req.Email),
	})
}

// @Router /auth/verify [post]
// @Summary Verify the user after sign up with code and get access_token.
// @Description Verify the user after sign up with code and get access_token.
// @Tags auth
// @Accept json
// @Produce json
// @Param data body models.VerifyRequest true "Data"
// @Success 200 {object} models.AuthResponse
// @Failure 500 {object} models.ResponseError
// @Failure 400 {object} models.ResponseError
// @Failure 404 {object} models.ResponseError
// @Failure 403 {object} models.ResponseError
func (h *handlerV1) Verify(ctx *gin.Context) {
	var (
		req models.VerifyRequest
	)

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	data, err := h.storage.User().GetByEmail(ctx.Request.Context(), req.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	if data != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(ErrEmailExists))
		return
	}

	userData, err := h.inMemory.Get("sign_up_" + req.Email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errResponse(ErrCodeExpired))
		return
	}

	var redisData models.AuthRequest
	err = json.Unmarshal([]byte(userData), &redisData)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errResponse(err))
		return
	}

	code, err := h.inMemory.Get("code_" + redisData.Email)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errResponse(ErrCodeExpired))
		return
	} else if code != req.Code {
		ctx.JSON(http.StatusForbidden, errResponse(ErrIncorrectCode))
		return
	}

	if code != req.Code {
		ctx.JSON(http.StatusForbidden, errResponse(ErrIncorrectCode))
		return
	}

	result, err := h.storage.User().CreateUser(ctx.Request.Context(), &repo.User{
		Email:    redisData.Email,
		Password: redisData.Password,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	token, _, err := utils.CreateToken(h.cfg, &utils.TokenParams{
		UserID:   result.ID,
		Email:    result.Email,
		Duration: time.Hour * 24,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.AuthResponse{
		AccessToken: token,
		ID:          result.ID,
		Email:       result.Email,
		CreatedAt:   result.CreatedAt,
	})
}

// @Router /auth/login [post]
// @Summary Login User
// @Description Login User
// @Tags auth
// @Accept json
// @Produce json
// @Param login body models.AuthRequest true "Login"
// @Success 200 {object} models.AuthResponse
// @Failure 500 {object} models.ResponseError
// @Failure 400 {object} models.ResponseError
// @Failure 404 {object} models.ResponseError
// @Failure 403 {object} models.ResponseError
func (h *handlerV1) Login(ctx *gin.Context) {
	var (
		req models.AuthRequest
	)

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	user, err := h.storage.User().GetByEmail(ctx.Request.Context(), req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errResponse(ErrEmailNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	err = utils.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errResponse(ErrWrongEmailOrPassword))
		return
	}

	token, _, err := utils.CreateToken(h.cfg, &utils.TokenParams{
		UserID:   user.ID,
		Email:    user.Email,
		Duration: time.Hour * 24,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.AuthResponse{
		ID:          user.ID,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		AccessToken: token,
	})
}
