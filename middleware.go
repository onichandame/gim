package gim

import "github.com/gin-gonic/gin"

type Middleware struct {
	Use gin.HandlerFunc
}
