package controller

import (
	"JD/dao"
	"JD/models"
	"JD/service"
	"JD/utils"
	_ "fmt"
	"github.com/gin-gonic/gin"
)

func AllShop(c *gin.Context) {
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
	slice := dao.AllShops()
	if slice == nil {
		service.ErrorReturn(c, nil)
		return
	}
	service.NormalReturn(c, *slice)
	return
}

func Chart(c *gin.Context) {

	var chart models.AddChart
	err := c.ShouldBind(&chart)
	if err != nil {
		service.ErrorReturn(c, "参数绑定失败")

		return
	}
	Info, exist := c.Get("Info")
	if !exist {
		service.ErrorReturn(c, "参数绑定失败")
	}
	BasicInfo, err := utils.Transform(Info)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}

	chart.BasicInfo = BasicInfo
	//info := <-s
	msg, err := dao.AddChart(chart)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	service.NormalReturn(c, msg)
	return

}

//c := <-m
//wait.Add(1)

func Update(c *gin.Context) {
	var chart models.ShopChart

	err := c.ShouldBind(&chart)
	if err != nil {
		service.ErrorReturn(c, "参数绑定失败")
		return
	}
	Info, exist := c.Get("Info")
	if !exist {
		service.ErrorReturn(c, "参数绑定失败")
		return
	}
	BasicInfo, err := utils.Transform(Info)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	chart.BasicInfo = BasicInfo
	//这里后面改一下
	ok, state := dao.UpdateChart(chart)
	c.JSON(200, gin.H{
		"state": ok,
		"msg":   state,
	})
	return

}

func AllChart(c *gin.Context) {
	var UserInfo models.Userinfo
	Info, exist := c.Get("Info")
	if !exist {
		service.ErrorReturn(c, "参数缺失")
	}
	BasicInfo, err := utils.Transform(Info)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	UserInfo.BasicInfo = BasicInfo
	msg, err := dao.AllChart(UserInfo)
	if err != nil {
		service.ErrorReturn(c, err.Error())
	}
	service.NormalReturn(c, *msg)
	return

}

func MakeOrder(c *gin.Context) {
	var Oder models.Order
	err := c.ShouldBind(&Oder)
	if err != nil {
		service.ErrorReturn(c, "参数绑定失败")
		return
	}
	//Oder.ChartId = c.PostFormArray("chart_id")
	Info, exist := c.Get("Info")
	if !exist {
		service.ErrorReturn(c, "参数绑定失败")
	}
	BasicInfo, err := utils.Transform(Info)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	Oder.BasicInfo = BasicInfo
	ok, state := dao.MakeOrder(Oder)
	c.JSON(200, gin.H{
		"state": ok,
		"msg":   state,
	})
	return
}

func GetCommit(c *gin.Context) {
	var commit models.Commits
	commit.Gid = c.Query("gid")
	Info, exist := c.Get("Info")
	if !exist {
		service.ErrorReturn(c, "参数缺失")
	}
	BasicInfo, err := utils.Transform(Info)
	if err != nil {
		service.ErrorReturn(c, err.Error())
		return
	}
	commit.BasicInfo = BasicInfo
	AllCommit := dao.GetCommit(commit)
	if AllCommit != nil {
		service.NormalReturn(c, *AllCommit)
		return
	}
	service.ErrorReturn(c, "操作失败")
}
func Show(c *gin.Context) {
	msg, err := dao.Class()
	if err != nil {
		c.JSON(200, gin.H{
			"state": false,
			"msg":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"state": true,
		"msg":   *msg,
	})
}
func Subscribe(c *gin.Context) {
	category := c.PostForm("gid")
	Info, exist := c.Get("Info")
	if !exist {
		service.ErrorReturn(c, "参数缺失")
	}
	BasicInfo, err := utils.Transform(Info)
	if err != nil {
		service.ErrorReturn(c, err.Error())
	}
	ok := utils.Subscribe(BasicInfo.Uid, category)
	if !ok {
		service.ErrorReturn(c, "订阅失败")
		return
	}
	service.NormalReturn(c, "订阅成功")
	return
}
