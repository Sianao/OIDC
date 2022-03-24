package hander

import (
	"JD/service"
	"JD/utils"
	"github.com/gin-gonic/gin"
)

func RootAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		Info, exist := c.Get("Info")
		if !exist {
			service.ErrorReturn(c, "参数绑定失败")
			c.Abort()
			return
		}
		BasicInfo, err := utils.Transform(Info)
		if err != nil {
			service.ErrorReturn(c, err.Error())
			c.Abort()
			return
		}
		if BasicInfo.Uid != 0 {
			service.ErrorReturn(c, "权限不够")
			c.Abort()
			return
		}
		c.Next()
	}

}
func UserConform() {

}
