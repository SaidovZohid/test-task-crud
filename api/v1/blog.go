package v1

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/SaidovZohid/test-task-crud/api/models"
	"github.com/SaidovZohid/test-task-crud/storage/repo"
	"github.com/gin-gonic/gin"
)

// @Security ApiKeyAuth
// @Router /blogs [post]
// @Summary Create a blog
// @Description Create a blog
// @Tags blog
// @Accept json
// @Produce json
// @Param blog body models.CreateAndUpdateBlog true "Blog"
// @Success 201 {object} models.Blog
// @Failure 500 {object} models.ResponseError
// @Failure 400 {object} models.ResponseError
func (h *handlerV1) CreateBlog(ctx *gin.Context) {
	var (
		req models.CreateAndUpdateBlog
	)

	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	userInfo, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	post, err := h.storage.Blog().CreateBlog(ctx.Request.Context(), &repo.Blog{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userInfo.UserID,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, models.Blog{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		User: &models.BlogUser{
			ID:    userInfo.UserID,
			Email: userInfo.Email,
		},
		CreatedAt: post.CreatedAt,
	})
}

// @Router /blogs/{id} [get]
// @Summary Get blog by its id
// @Description Get blog by its id
// @Tags blog
// @Produce json
// @Param id path int true "Blog ID"
// @Success 200 {object} models.Blog
// @Failure 500 {object} models.ResponseError
// @Failure 404 {object} models.ResponseError
func (h *handlerV1) GetBlogByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	blog, err := h.storage.Blog().GetBlog(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errResponse(ErrNoBlogFound))
		} else {
			ctx.JSON(http.StatusInternalServerError, errResponse(err))
		}
		return
	}
	user, err := h.storage.User().GetByUserID(ctx.Request.Context(), blog.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.Blog{
		ID:      id,
		Title:   blog.Title,
		Content: blog.Content,
		User: &models.BlogUser{
			Email: user.Email,
			ID:    user.ID,
		},
		CreatedAt: blog.CreatedAt,
	})
}

// @Security ApiKeyAuth
// @Router /blogs/{id} [put]
// @Summary Update blog with it's id as param
// @Description Update blog with it's id as param
// @Tags blog
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param post body models.CreateAndUpdateBlog true "Post"
// @Success 200 {object} models.Blog
// @Failure 500 {object} models.ResponseError
// @Failure 400 {object} models.ResponseError
// @Failure 404 {object} models.ResponseError
func (h *handlerV1) UpdateBlog(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	var (
		req models.CreateAndUpdateBlog
	)

	err = ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	userInfo, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	blog, err := h.storage.Blog().Update(ctx.Request.Context(), &repo.Blog{
		ID:      id,
		Title:   req.Title,
		Content: req.Content,
		UserID:  userInfo.UserID,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errResponse(ErrNoBlogFound))
		} else {
			ctx.JSON(http.StatusInternalServerError, errResponse(err))
		}
		return
	}

	ctx.JSON(http.StatusOK, models.Blog{
		ID:      blog.ID,
		Title:   blog.Title,
		Content: blog.Content,
		User: &models.BlogUser{
			ID:    userInfo.UserID,
			Email: userInfo.Email,
		},
		CreatedAt: blog.CreatedAt,
	})
}

// @Security ApiKeyAuth
// @Router /blogs/{id} [delete]
// @Summary Delete blog
// @Description Delete blog works only for user who created it!
// @Tags blog
// @Produce json
// @Param id path int true "Blog ID"
// @Success 200 {object} models.Message
// @Failure 500 {object} models.ResponseError
// @Failure 404 {object} models.ResponseError
// @Failure 400 {object} models.ResponseError
func (h *handlerV1) DeleteBlog(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	userInfo, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	err = h.storage.Blog().DeleteBlog(ctx.Request.Context(), id, userInfo.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errResponse(ErrNoBlogFound))
		} else {
			ctx.JSON(http.StatusInternalServerError, errResponse(err))
		}
		return
	}

	ctx.JSON(http.StatusOK, models.Message{
		Msg: "Successfully deleted!",
	})
}

// @Router /blogs [get]
// @Summary Get blogs by giving limit, page, user_id, sort and search
// @Description Get blogs by giving limit, page, user_id, sort and search
// @Tags blog
// @Produce json
// @Param filter query models.GetAllParams false "Filter"
// @Success 200 {object} models.GetAllBlogsResult
// @Failure 500 {object} models.ResponseError
// @Failure 400 {object} models.ResponseError
func (h *handlerV1) GetAllPosts(ctx *gin.Context) {
	params, err := validateGetAllParams(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	result, err := h.storage.Blog().GetAll(ctx.Request.Context(), &repo.GetBlogsParams{
		Limit:      params.Limit,
		Page:       params.Page,
		Search:     params.Search,
		UserID:     params.UserID,
		SortByDate: params.SortByDate,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, getPostsResponse(result))
}

func getPostsResponse(data *repo.GetAllBlogsResult) *models.GetAllBlogsResult {
	response := models.GetAllBlogsResult{
		Blogs: make([]*models.Blog, 0),
		Count: data.Count,
	}

	for _, post := range data.Blogs {
		u := parsePostModel(post)
		response.Blogs = append(response.Blogs, &u)
	}

	return &response
}

func parsePostModel(post *repo.Blog) models.Blog {
	return models.Blog{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		User: &models.BlogUser{
			ID:    post.UserID,
			Email: "",
		},
		CreatedAt: post.CreatedAt,
	}
}
