package handles

import (
	"midea-clo/db"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type ApprovalParams struct {
	Mobile  string `form:"mobile" binding:"required"`
	IsAgree bool   `form:"isagree"`
}

func Approval(c *gin.Context) {
	var params ApprovalParams
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

	if params.IsAgree {
		info.Status = 1
	} else {
		info.Status = 2
	}

	info.AgreeTime = db.NullTime{mysql.NullTime{time.Now(), true}}
	result := db.MysqlDB.Save(&info).RowsAffected

	if result == 0 {
		Response(c, CODE_FAIL, "审批失败", nil)
		return
	}

	Response(c, CODE_OK, "审批成功", nil)
}
