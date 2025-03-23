package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/varsilias/simplebank/token"
)

const (
	authorisationHeaderKey  = "authorization"
	authorisationTypeBearer = "bearer"
	authorisationKey        = "user"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authorisationHeaderKey)
		if len(authHeader) == 0 {
			err := errors.New("please provide a valid access token")
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(http.StatusForbidden, ctx.Request.URL.Path, err))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(http.StatusForbidden, ctx.Request.URL.Path, err))
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != authorisationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authType)
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(http.StatusForbidden, ctx.Request.URL.Path, err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(http.StatusUnauthorized, ctx.Request.URL.Path, err))
			return

		}

		ctx.Set(authorisationKey, payload)
		ctx.Next()
	}
}
