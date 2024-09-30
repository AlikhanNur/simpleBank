package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/alikhanMuslim/simpleBank/token"
	"github.com/gin-gonic/gin"
)

const (
	authHeaderkey  = "authorization"
	authTypeBearer = "bearer"
	authPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.PasetoMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authHeaderkey)
		if len(authHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != authTypeBearer {
			err := fmt.Errorf("unsuported authorization type %s", authType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accesstoken := fields[1]
		payload, err := tokenMaker.VerifyToken(accesstoken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authPayloadKey, payload)
		ctx.Next()
	}
}
