package controller

import (
	"JD/dao"
	"JD/models"
	"JD/service"
	"JD/utils"
	"github.com/gin-gonic/gin"
)

func SendMessage(c *gin.Context) {
	number := c.Query("Phone")
	//调用工具发送短信
	err := utils.Sendsms(number)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	service.NormalReturn(c, "短信发送成功")
	return
}

func Register(c *gin.Context) {
	//接受参数
	var register models.Register

	err := c.ShouldBind(&register)
	if err != nil {
		service.ErrorReturn(c, "参数绑定失败")
		return
	}
	//参数传递 进行校验
	err = utils.GetSmsConform(register.Number, register.Code)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	err = dao.Register(register.Username, register.Password, register.Number)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	service.NormalReturn(c, "注册成功")
	return

}
