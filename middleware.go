package gim

import "github.com/gin-gonic/gin"

type withMiddleware interface {
	Use() gin.HandlerFunc
}
