package handles

import (
	"crypto/sha1"
	"fmt"
	"midea-clo/util"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

type SignatureParams struct {
	Signature string `form:"signature"`
	Timestamp string `form:"timestamp"`
	Nonce     string `form:"nonce"`
	Echostr   string `form:"echostr"`
}

func CheckSignature(c *gin.Context) {
	var params SignatureParams
	// 验证数据并绑定
	if err := c.ShouldBind(&params); err != nil {
		ResponseError(c, err)
		return
	}

	signature := []string{util.WechatManger.Token, params.Timestamp, params.Nonce}
	sort.Strings(signature)
	h := sha1.New()
	h.Write([]byte(strings.Join(signature, "")))
	bs := h.Sum(nil)

	if fmt.Sprintf("%x", bs) != params.Signature {
		c.String(http.StatusOK, "false")
		return
	}

	c.String(http.StatusOK, params.Echostr)
}
