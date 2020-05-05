package handles

import (
	"midea-clo/db"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type createSignUpParams struct {
	Name     string `form:"name" binding:"required"`
	Company  string `form:"company" binding:"required"`
	Position string `form:"position" binding:"required"`
	Mobile   string `form:"mobile" binding:"required"`
	// Code     string `form:"code" binding:"required"`
}

func SignUp(c *gin.Context) {
	var params createSignUpParams
	// 验证数据并绑定
	if err := c.ShouldBindJSON(&params); err != nil {
		ResponseError(c, err)
		return
	}

	if ok, _ := regexp.MatchString(`^(1[0-9]{10})$`, params.Mobile); !ok {
		Response(c, CODE_PARAMS_ERR, "手机号码格式错误", nil)
		return
	}

	var info db.Info
	if !db.MysqlDB.Where("mobile = ?", params.Mobile).First(&info).RecordNotFound() {
		Response(c, CODE_SIGNUP_REPEAT, "您已报名成功，无需重复报名", nil)
		return
	}

	var createInfo = db.Info{
		Name:       params.Name,
		Company:    params.Company,
		Position:   params.Position,
		Mobile:     params.Mobile,
		CreateTime: db.NullTime{mysql.NullTime{time.Now(), true}},
	}

	result := db.MysqlDB.Create(&createInfo).RowsAffected
	if result == 0 {
		Response(c, CODE_FAIL, "报名失败", nil)
		return
	}

	Response(c, CODE_OK, "报名成功", nil)
}