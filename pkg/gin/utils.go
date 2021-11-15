package gin

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	goutils "github.com/onichandame/go-utils"
)

type withStatus interface {
	Status() int
}
type withBody interface {
	Body() interface{}
}

func GetHTTPHandler(fn func(*gin.Context) interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		var res interface{}
		err := goutils.Try(func() { res = fn(c) })
		var status int
		var contentType string
		var response []byte
		var body interface{}
		populateRes := func(res interface{}, defStatus int, defBody interface{}) {
			if s, ok := res.(withStatus); ok {
				status = s.Status()
			} else {
				status = defStatus
			}
			if r, ok := res.(withBody); ok {
				body = r.Body()
			} else {
				body = defBody
			}
		}
		if err == nil {
			populateRes(res, 200, res)
		} else {
			populateRes(err, 400, err.Error())
		}
		if b, ok := body.([]byte); ok {
			response = b
			contentType = "text/plain"
		} else if r, ok := body.(string); ok {
			response = []byte(r)
			contentType = "text/plain"
		} else {
			if response, err = json.Marshal(body); err != nil {
				status = 500
				response = []byte("failed to serialize response body")
				contentType = "text/plain"
				fmt.Println(err)
			} else {
				contentType = "application/json"
			}
		}
		c.Data(status, contentType, response)
	}
}
