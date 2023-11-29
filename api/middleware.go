package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"simplebank/token"
	"strings"
)

const authorizationHeaderKey = "authorization"
const authorizationTypeBearer = "bearer"
const autuorzationPayloadKey = "authorization_key"

func authMiddleware(token token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not previded")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		authorizationHeaderType := strings.ToLower(fields[0])
		if authorizationHeaderType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationTypeBearer)
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := token.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		c.Set(autuorzationPayloadKey, *payload)
		c.Next()
	}
}
