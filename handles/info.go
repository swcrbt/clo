package handles

import (
	"midea-clo/db"

	"github.com/gin-gonic/gin"
)

type GetInfoParams struct {
	Mobile string `form:"mobile"`
}

func Info(c *gin.Context) {
	var params GetInfoParams
	// 验证数据并绑定
	if err := c.ShouldBind(&params); err != nil {
		ResponseError(c, err)
		return
	}

	var info db.Info
	if db.MysqlDB.Where("mobile = ?", params.Mobile).First(&info).RecordNotFound() {
		Response(c, CODE_NOT_RECORD, "信息不存在", nil)
		return
	}

	Response(c, CODE_OK, "获取成功", info)
}
