package v1

import (
	"errors"
	"strconv"

	"github.com/SaidovZohid/test-task-crud/api/models"
	"github.com/SaidovZohid/test-task-crud/config"
	"github.com/SaidovZohid/test-task-crud/pkg/logger"
	"github.com/SaidovZohid/test-task-crud/storage"
	"github.com/gin-gonic/gin"
)

var (
	ErrWrongEmailOrPassword = errors.New("wrong email or password")
	ErrUserNotVerifid       = errors.New("user not verified")
	ErrEmailExists          = errors.New("email is already exists")
	ErrIncorrectCode        = errors.New("incorrect verification code")
	ErrEmailNotFound        = errors.New("email does not exist")
	ErrCodeExpired          = errors.New("verification code is expired")
	ErrNoBlogFound          = errors.New("blog does not exist")
	ErrForbidden            = errors.New("forbidden")
)

type handlerV1 struct {
	cfg      *config.Config
	storage  storage.StorageI
	inMemory storage.InMemoryStorageI
	log      logger.Logger
}

type HandlerV1Options struct {
	Cfg      *config.Config
	Storage  storage.StorageI
	InMemory storage.InMemoryStorageI
	Log      logger.Logger
}

func New(options *HandlerV1Options) *handlerV1 {
	return &handlerV1{
		cfg:      options.Cfg,
		storage:  options.Storage,
		inMemory: options.InMemory,
		log:      options.Log,
	}
}

func errResponse(err error) *models.ResponseError {
	return &models.ResponseError{
		Error: err.Error(),
	}
}

func validateGetAllParams(ctx *gin.Context) (*models.GetAllParams, error) {
	var (
		limit      int64 = 10
		page       int64 = 1
		err        error
		userId     int64
		sortByDate string
	)
	if ctx.Query("limit") != "" {
		limit, err = strconv.ParseInt(ctx.Query("limit"), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	if ctx.Query("page") != "" {
		page, err = strconv.ParseInt(ctx.Query("page"), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	if ctx.Query("user_id") != "" {
		userId, err = strconv.ParseInt(ctx.Query("user_id"), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	if ctx.Query("sort") != "" &&
		(ctx.Query("sort") == "desc" || ctx.Query("sort") == "asc") {
		sortByDate = ctx.Query("sort")
	}

	return &models.GetAllParams{
		Limit:      limit,
		Page:       page,
		Search:     ctx.Query("search"),
		UserID:     userId,
		SortByDate: sortByDate,
	}, nil
}
