package handles

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CODE_OK = iota // Success
	CODE_FAIL
	CODE_PARAMS_ERR
	CODE_PARAMS_MISS
	CODE_SIGNUP_REPEAT
	CODE_NOT_RECORD
)

func Response(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

func ResponseError(c *gin.Context, err error) {
	fmt.Println(fmt.Sprintf("%T", err))
	fmt.Println(err)
	if fmt.Sprintf("%T", err) == "validator.ValidationErrors" {
		Response(c, CODE_PARAMS_ERR, "缺少参数", nil)
	} else {
		Response(c, CODE_PARAMS_MISS, "参数错误", nil)
	}
}
