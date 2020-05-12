package handles

import (
	"midea-clo/db"

	"github.com/gin-gonic/gin"
)

type GetListParams struct {
	Page     int `form:"page"`
	PageSize int `form:"pagesize"`
	Status   int `form:"status"`
}

func List(c *gin.Context) {
	var params GetListParams
	// 验证数据并绑定
	if err := c.ShouldBind(&params); err != nil {
		ResponseError(c, err)
		return
	}

	var infos []db.Info
	var count int
	listQuery := db.MysqlDB.Order("createtime desc")
	countQuery := db.MysqlDB.Model(&db.Info{})

	if params.Status >= 0 {
		listQuery = listQuery.Where("status = ?", params.Status)
		countQuery = countQuery.Where("status = ?", params.Status)
	}

	if params.PageSize > 0 {
		if params.Page > 1 {
			listQuery = listQuery.Offset((params.Page - 1) * params.PageSize)
		}

		listQuery = listQuery.Limit(params.PageSize)
	}

	if listQuery.Find(&infos).RecordNotFound() {
		Response(c, CODE_NOT_RECORD, "列表不存在", nil)
		return
	}

	countQuery.Count(&count)

	Response(c, CODE_OK, "获取成功", gin.H{
		"page":        params.Page,
		"pagesize":    params.PageSize,
		"recordcount": count,
		"list":        infos,
	})
}
