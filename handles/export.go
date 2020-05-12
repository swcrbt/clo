package handles

import (
	"bytes"
	"fmt"
	"midea-clo/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
)

type ExportParams struct {
	Filename string `form:"filename"`
}

func Export(c *gin.Context) {
	var params ExportParams
	// 验证数据并绑定
	if err := c.ShouldBind(&params); err != nil {
		ResponseError(c, err)
		return
	}

	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("sheet")
	// 设置表格头
	row := sheet.AddRow()
	var headers = []string{"姓名", "联系电话", "公司", "职位", "省份", "关系", "公司性质", "应用领域", "审批结果", "是否签到"}
	for _, header := range headers {
		row.AddCell().Value = header
	}
	var infos []db.Info
	db.MysqlDB.Order("createtime desc").Find(&infos)
	// 写入数据
	for _, info := range infos {
		row := sheet.AddRow()
		row.AddCell().Value = info.Name
		row.AddCell().Value = info.Mobile
		row.AddCell().Value = info.Company
		row.AddCell().Value = info.Position
		row.AddCell().Value = info.Province

		row.AddCell().Value = db.RelationText[info.Relation]
		row.AddCell().Value = info.Nature
		row.AddCell().Value = info.Region

		if info.Status == 1 {
			row.AddCell().Value = "通过"
		} else if info.Status == 2 {
			row.AddCell().Value = "不通过"
		} else {
			row.AddCell().Value = "未审批"
		}

		if info.IsSign {
			row.AddCell().Value = "已签到"
		} else {
			row.AddCell().Value = "未签到"
		}
	}
	// c.Header("Content-Disposition", "attachment; filename="+fmt.Sprintf("%s", "file.xls")) //文件名
	// c.Header("Content-Disposition", "attachment")
	// c.Header("Content-Type", "application/vnd.ms-excel")
	// xlsx
	// c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	fileName := "file.xlsx"
	if params.Filename != "" {
		fileName = params.Filename
	}
	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename="` + fileName + `"`,
	}

	var buffer bytes.Buffer
	if err := file.Write(&buffer); err != nil {
		fmt.Println(err)
	}
	r := bytes.NewReader(buffer.Bytes())
	c.DataFromReader(
		http.StatusOK,
		int64(r.Len()),
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		r,
		extraHeaders,
	)
}
