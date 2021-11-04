package graphql

import (
	"context"

	"github.com/gin-gonic/gin"
)

func GetContext(c context.Context) *gin.Context {
	return c.Value(contextToken).(*gin.Context)
}
