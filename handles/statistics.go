package handles

import (
	"midea-clo/db"

	"github.com/gin-gonic/gin"
)

func Statistics(c *gin.Context) {
	var total int
	var sign int

	db.MysqlDB.Model(&db.Info{}).Count(&total)
	db.MysqlDB.Model(&db.Info{}).Where("issign = ?", "1").Count(&sign)

	Response(c, CODE_OK, "获取成功", gin.H{
		"total": total,
		"sign":  sign,
	})
}
