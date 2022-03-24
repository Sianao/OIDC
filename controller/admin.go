package controller

import (
	"JD/dao"
	"JD/models"
	"JD/service"
	"JD/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"time"
)

func RootLogin(c *gin.Context) {
	var admin models.Admin
	err := c.ShouldBind(&admin)
	if err != nil {
		service.ErrorReturn(c, "参数绑定失败")
	}
	ok, state := dao.AdminLogin(admin.Name, admin.Password)

	if !ok {
		service.ErrorReturn(c, state)
		return
	}
	var BasicInfo models.BasicInfo
	BasicInfo.Uid = 0
	BasicInfo.Username = admin.Name
	token := utils.MakeToken(BasicInfo)
	ok = utils.SetToken(1, token)
	if !ok {
		service.ErrorReturn(c, "登录失败")
		return

	}
	service.NormalReturn(c, gin.H{
		"code":  ok,
		"msg":   state,
		"token": token,
	})
	return

}
func RootLogout(c *gin.Context) {
	Authorization := c.Request.Header.Get("Authorization")
	ok := utils.DeleteToken(Authorization)
	if !ok {
		service.ErrorReturn(c, "推出登录失败")
		return
	}
	service.NormalReturn(c, "退出登录成功")
	return
}
func RootAll(c *gin.Context) {
	Allorder, err := dao.GetAllOrder()
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	service.NormalReturn(c, *Allorder)
	return
}
func UpdateONeOrder(c *gin.Context) {
	var update models.UpdateUserOrder
	err := c.ShouldBind(&update)
	if err != nil {
		service.ErrorReturn(c, "参数绑定错误")
		return
	}
	msg, err := dao.OrderChange(update)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	service.NormalReturn(c, msg)
	return
}
func DeleteUserOrder(c *gin.Context) {
	var update models.UpdateUserOrder
	err := c.ShouldBind(&update)
	if err != nil {
		service.ErrorReturn(c, "参数绑定失败")
	}
	msg, err := dao.DeleteUserOrder(update)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	service.NormalReturn(c, msg)
	return
}
func AddGoods(c *gin.Context) {
	var add models.GoodsAdd
	err := c.ShouldBind(&add)
	if err != nil {

		err = errors.New("参数绑定失败")
		service.ErrorReturn(c, err.Error())
		return
	}
	t := time.Now()
	utils.AddNewInfo(add.Category, "商品发布", "您关注的分类发布了新的商品", t)
	msg, err := dao.AddGoods(add, c)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	service.NormalReturn(c, msg)
	return
}

func UpdateGoods(c *gin.Context) {
	var update models.UpdateGoods
	err := c.ShouldBind(&update)
	if err != nil {
		err = errors.New("参数绑定失败")
		service.ErrorReturn(c, err.Error())
	}
	msg, err := dao.UpdateGoods(update, c)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	service.NormalReturn(c, msg)
	return
}
func DeleteGoods(c *gin.Context) {
	Gid := c.PostForm("Gid")
	msg, err := dao.DeleteGoods(Gid)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	service.NormalReturn(c, msg)
	return

}
func AddCategory(c *gin.Context) {
	category := c.PostForm("category")
	err := utils.AddCategory(category)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	service.NormalReturn(c, "添加成功")
	return
}
func AllCategory(c *gin.Context) {
	all := utils.AllCategory()
	service.NormalReturn(c, all)
	return

}
