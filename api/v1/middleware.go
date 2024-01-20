package v1

import (
	"errors"
	"net/http"

	"github.com/SaidovZohid/test-task-crud/pkg/utils"
	"github.com/gin-gonic/gin"
)

func (h *handlerV1) AuthMiddleWare(ctx *gin.Context) {
	accessToken := ctx.GetHeader(h.cfg.AuthHeaderKey)

	if len(accessToken) == 0 {
		err := errors.New("unauthenticated user")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
		return
	}
	payload, err := utils.VerifyToken(h.cfg, accessToken)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	ctx.Set(h.cfg.AuthPayloadKey, payload)
	ctx.Next()
}

func (h *handlerV1) GetAuthPayload(ctx *gin.Context) (*utils.Payload, error) {
	i, exist := ctx.Get(h.cfg.AuthPayloadKey)
	if !exist {
		return nil, errors.New("not found payload")
	}

	payload, ok := i.(*utils.Payload)
	if !ok {
		return nil, errors.New("unknown user")
	}
	return payload, nil
}
