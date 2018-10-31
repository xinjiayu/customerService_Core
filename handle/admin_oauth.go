package handle

import (
	"encoding/base64"
	"errors"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// 管理员后台鉴权
func AdminOauthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err   error
			kfId  string
			token = c.Request.Header.Get("authentication")
		)

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "API token required",
			})
			c.Abort()
			return
		}

		if kfId, err = adminAuthToken2Model(c); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  err.Error(),
			})
			c.Abort()
			return
		} else {
			c.Set("KFID", kfId)
		}

		// 更新在线客服列表时间
		contrastKf := model.Kf{Id: kfId, UpdateTime: time.Now()}
		isExist := contrastKf.OnlineExist()
		if isExist {
			model.KfLastTimeChange <- &contrastKf
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "授权过期，请重新登录",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// 鉴权Token解码为模型
func adminAuthToken2Model(c *gin.Context) (kfId string, err error) {
	var (
		token = c.Request.Header.Get("authentication")
		aes   = common.AesEncrypt{}
	)

	decodeBytes, err := base64.StdEncoding.DecodeString(token)
	if bytes, err := aes.Decrypt(decodeBytes); err != nil {
		err = errors.New("API token Authentication failed")
	} else {
		kfId = string(bytes)
		if kfId == "" {
			err = errors.New("API token Authentication failed")
		}
	}
	return
}