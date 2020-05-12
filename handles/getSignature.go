package handles

import (
	"crypto/sha1"
	"fmt"
	"midea-clo/util"

	"github.com/gin-gonic/gin"
)

type GetSignatureParams struct {
	Timestamp int    `form:"timestamp" binding:"required"`
	NonceStr  string `form:"noncestr" binding:"required"`
	Url       string `form:"url" binding:"required"`
}

// 巨坑：JS接口安全域名不需要协议
func GetSignature(c *gin.Context) {
	var params GetSignatureParams
	// 验证数据并绑定
	if err := c.ShouldBind(&params); err != nil {
		ResponseError(c, err)
		return
	}

	str := fmt.Sprintf(
		"jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s",
		util.WechatManger.GetApiTicket(),
		params.NonceStr,
		params.Timestamp,
		params.Url,
	)

	h := sha1.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)

	Response(c, CODE_OK, "获取成功", gin.H{
		"appid":     util.WechatManger.AppID,
		"noncestr":  params.NonceStr,
		"timestamp": params.Timestamp,
		"url":       params.Url,
		"signature": fmt.Sprintf("%x", bs),
	})
}
