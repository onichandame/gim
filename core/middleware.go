package core

import "github.com/gin-gonic/gin"

type Middleware struct {
	Use gin.HandlerFunc
}
