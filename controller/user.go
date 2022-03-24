package controller

import (
	"JD/dao"
	"JD/models"
	"JD/service"
	"JD/utils"

	"github.com/gin-gonic/gin"
)

func BalanceGet(c *gin.Context) {
	var user models.User
	Info, exist := c.Get("Info")
	if !exist {
		service.ErrorReturn(c, "参数缺失")

	}
	BasicInfo, err := utils.Transform(Info)
	if err != nil {
		service.ErrorReturn(c, "参数缺失")
		return
	}
	user.BasicInfo = BasicInfo
	balance, err := dao.GetBalance(user.Username)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	c.JSON(200, gin.H{
		"state":   true,
		"msg":     "查找成功",
		"balance": *balance,
	})

}
func BalanceCharge(c *gin.Context) {
	var user models.Balance
	err := c.ShouldBind(&user)
	if err != nil {
		service.ErrorReturn(c, "参数绑定失败")
		return
	}
	Info, exist := c.Get("Info")
	if !exist {
		service.ErrorReturn(c, "参数缺失")
	}
	BasicInfo, err := utils.Transform(Info)

	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	user.BasicInfo = BasicInfo

	ok, state := dao.ChargeBalance(user)
	if !ok {
		service.ErrorReturn(c, state)
		return
	}
	service.NormalReturn(c, state)
	return

}
func Order(c *gin.Context) {
	var user models.User
	Info, exist := c.Get("Info")
	if !exist {
		service.ErrorReturn(c, "参数缺失")
	}
	BasicInfo, err := utils.Transform(Info)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	user.BasicInfo = BasicInfo
	ok, info := dao.AllOrder(user)
	if ok {
		c.JSON(200, *info)
		return
	}
	c.JSON(200, gin.H{
		"state": "false",
	})
}
func UpdateOrder(c *gin.Context) {
	var order models.UpdateOrder
	err := c.ShouldBind(&order)
	if err != nil {
		service.ErrorReturn(c, "参数缺失")
		return
	}
	Info, exist := c.Get("Info")
	if !exist {
		service.ErrorReturn(c, "参数缺失")
	}
	BasicInfo, err := utils.Transform(Info)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	order.BasicInfo = BasicInfo
	err = dao.UpdateOrder(order)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	service.NormalReturn(c, "确认收货成功")
	return

}
func DeleteOrder(c *gin.Context) {
	var order models.UpdateOrder
	err := c.ShouldBind(&order)
	if err != nil {
		service.ErrorReturn(c, "参数缺失")
		return
	}
	Info, exist := c.Get("Info")
	if !exist {
		service.ErrorReturn(c, "参数缺失")
	}
	BasicInfo, err := utils.Transform(Info)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	order.BasicInfo = BasicInfo
	err = dao.DeleteOrder(order)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	service.NormalReturn(c, "订单销毁成功")
}
func Commit(c *gin.Context) {
	var Commit models.Commit
	err := c.ShouldBind(&Commit)
	if err != nil {
		service.ErrorReturn(c, "参数缺失")
		return
	}
	Info, exist := c.Get("Info")
	if !exist {
		service.ErrorReturn(c, "参数缺失")
	}
	BasicInfo, err := utils.Transform(Info)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	Commit.BasicInfo = BasicInfo
	ok, msg := dao.Commit(Commit)
	c.JSON(200, gin.H{
		"state": ok,
		"msg":   msg,
	})
}

// ImageUser 用户头像修改
func ImageUser(c *gin.Context) {
	var User models.UserImage
	err := c.ShouldBind(&User)
	if err != nil {
		service.ErrorReturn(c, "参数绑定失败")
		return
	}
	Info, exist := c.Get("Info")
	if !exist {
		service.ErrorReturn(c, "参数缺失")
	}
	BasicInfo, err := utils.Transform(Info)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	User.BasicInfo = BasicInfo
	url, err := utils.SaveFile(User.Image, c)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	msg, err := dao.SaveFile(url, User.BasicInfo)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	service.NormalReturn(c, msg)
	return
}
func Info(c *gin.Context) {
	var User models.BasicInfo
	//从中间件的set 中取出token的荷载信息 //
	//原本这玩意是放在body里 后面又说放在head里

	Info, exist := c.Get("Info")
	if !exist {
		service.ErrorReturn(c, "参数缺失")
	}
	BasicInfo, err := utils.Transform(Info)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	User = BasicInfo
	UserInfo, err := dao.MyInfo(User)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	service.NormalReturn(c, *UserInfo)
	return
}
func GetInfo(c *gin.Context) {
	Info, exist := c.Get("Info")
	if !exist {
		service.ErrorReturn(c, "参数缺失")
	}
	BasicInfo, err := utils.Transform(Info)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	info := utils.GetInfo(BasicInfo.Uid)
	if info != nil {
		service.NormalReturn(c, info)
		return
	}
	service.ErrorReturn(c, "查询错误")
	return

}
func UserCategory(c *gin.Context) {
	Info, exist := c.Get("Info")
	if !exist {
		service.ErrorReturn(c, "参数缺失")
	}
	BasicInfo, err := utils.Transform(Info)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	all := utils.MySubscribe(BasicInfo.Uid)
	service.NormalReturn(c, all)
	return
}
func UnSubscribe(c *gin.Context) {
	todo := c.PostForm("category")
	Info, exist := c.Get("Info")
	if !exist {
		service.ErrorReturn(c, "参数缺失")
	}
	BasicInfo, err := utils.Transform(Info)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	err = utils.Unsubscribe(BasicInfo.Uid, todo)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	service.NormalReturn(c, "退订成功")
	return

}
