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
	Province string `form:"province"`
	Relation int    `form:"relation"`
	Nature   string `form:"nature"`
	Region   string `form:"region"`
	// Code     string `form:"code" binding:"required"`
}

func SignUp(c *gin.Context) {
	var params createSignUpParams
	// 验证数据并绑定
	if err := c.ShouldBind(&params); err != nil {
		ResponseError(c, err)
		return
	}

	if ok, _ := regexp.MatchString(`^(1[0-9]{10})$`, params.Mobile); !ok {
		Response(c, CODE_PARAMS_ERR, "手机号码格式错误", nil)
		return
	}

	var info db.Info
	if !db.MysqlDB.Where("mobile = ?", params.Mobile).First(&info).RecordNotFound() {
		if info.Status == 1 {
			Response(c, CODE_SIGNUP_REPEAT_AGREE, "您已报名成功，审批通过", nil)
			return
		}

		if info.Status == 2 {
			Response(c, CODE_SIGNUP_REPEAT_UNAGREE, "您已报名成功，审批不通过", nil)
			return
		}

		Response(c, CODE_SIGNUP_REPEAT, "您已报名成功，无需重复报名", nil)
		return
	}

	var createInfo = db.Info{
		Name:       params.Name,
		Company:    params.Company,
		Position:   params.Position,
		Mobile:     params.Mobile,
		Province:   params.Province,
		Relation:   params.Relation,
		Nature:     params.Nature,
		Region:     params.Region,
		CreateTime: db.NullTime{mysql.NullTime{time.Now(), true}},
	}

	result := db.MysqlDB.Create(&createInfo).RowsAffected
	if result == 0 {
		Response(c, CODE_FAIL, "报名失败", nil)
		return
	}

	Response(c, CODE_OK, "报名成功", nil)
}
