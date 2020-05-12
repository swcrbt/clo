package handles

import (
	"midea-clo/db"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type SignParams struct {
	Mobile string `form:"mobile" binding:"required"`
}

func Sign(c *gin.Context) {
	var params SignParams
	// 验证数据并绑定
	if err := c.ShouldBind(&params); err != nil {
		ResponseError(c, err)
		return
	}

	var info db.Info
	if db.MysqlDB.Where("mobile = ?", params.Mobile).First(&info).RecordNotFound() {
		Response(c, CODE_NOT_RECORD, "您还没有报名", nil)
		return
	}

	if info.Status == 0 {
		Response(c, CODE_NOT_APPROVED, "您的报名未审批", nil)
		return
	}

	if info.Status != 1 {
		Response(c, CODE_NOT_AGREE, "您的报名不通过", nil)
		return
	}

	info.IsSign = true
	info.SignTime = db.NullTime{mysql.NullTime{time.Now(), true}}
	result := db.MysqlDB.Save(&info).RowsAffected

	if result == 0 {
		Response(c, CODE_FAIL, "签到失败", nil)
		return
	}

	Response(c, CODE_OK, "签到成功", info)
}
