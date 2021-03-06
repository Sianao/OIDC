package router

import (
	"JD/controller"
	"JD/hander"
	"github.com/gin-gonic/gin"
)

func Entrance() {

	r := gin.Default()
	//使用中间件 获取用户部分状态
	r.Use(hander.Cors())
	//静态文件加载 但是貌似 有一丢丢慢
	//client 本来是准备用来写个client端发请求的
	//
	//
	// OIDC 实现主要是靠下面五个接口实现的  这里简述一下流程实现
	//
	// 首先 用户携带 client_id scope response_typed state 等参数访问 oauth_login 接口
	r.GET("/client", controller.Client)
	//接收数据 进行保存  跳转到登录接口 并携带必要参数 部分数据存入 redis
	r.GET("/oauth_login", controller.Service)
	//Login 接口 进行数据接收 并发送  临时生成code  和token 取出 之前的访问参数 携带state code 返回
	r.POST("/login", controller.Login)
	// 用户携带授权码  申请获取token 并对用户信息进行检验 返回用户token
	r.GET("/oauth", controller.Oauth)

	//  不足 : 没有对参数进行严格的检验 只相当于走了个流程 其中的一些跳转 数据保存问题处理欠缺
	//
	//
	//

	r.GET("/callback", controller.Callback)

	//
	///
	//注册短信发送接口
	r.GET("/register", controller.SendMessage)
	//注册接口
	r.POST("/register", controller.Register)
	//登录接口

	//找回密码
	r.POST("/find", controller.Find)
	//登出接口
	r.GET("/logout", controller.Logout)
	r.GET("/show", controller.Show)
	r.GET("/news", controller.News)
	//用户组
	UserGroup := r.Group("/user")
	{
		//用户主界面
		UserGroup.GET("/", hander.Auth(), controller.Info)
		//更新用户信息
		UserGroup.PUT("/image", hander.Auth(), controller.ImageUser)
		//Post用于充值
		UserGroup.POST("/balance", hander.Auth(), controller.BalanceCharge)
		//GET用于查询
		UserGroup.GET("/balance", hander.Auth(), controller.BalanceGet)
		//用于获取用户订单
		UserGroup.GET("/order", hander.Auth(), controller.Order)
		//更新用户订单
		UserGroup.PUT("/order", hander.Auth(), controller.UpdateOrder)
		//删除订单
		UserGroup.DELETE("/order", hander.Auth(), controller.DeleteOrder)
		//Post 提交评论
		UserGroup.POST("/commit", hander.Auth(), controller.Commit)
		//获取用户订阅信息
		UserGroup.GET("/subscribe", hander.Auth(), controller.UserCategory)
		//退订
		UserGroup.DELETE("/subscribe", hander.Auth(), controller.UnSubscribe)
		//用户获取信息
		UserGroup.GET("/message", hander.Auth(), controller.GetInfo)
		UserGroup.POST("/readed", hander.Auth(), controller.MarkReaded)
	}
	ShopCenter := r.Group("/shop")
	{
		//all 显示所有商品
		ShopCenter.GET("/all", hander.Auth(), controller.AllShop)
		//订阅评道
		ShopCenter.POST("/subscribe", hander.Auth(), controller.Subscribe)
		//commit 获取评论
		ShopCenter.GET("/commit", hander.Auth(), controller.GetCommit)
		//post 添加商品
		ShopCenter.POST("/chart", hander.Auth(), controller.Chart)
		//Get 获取购物车信息
		ShopCenter.GET("/chart", hander.Auth(), controller.AllChart)
		//update 对购物车信息进行修改
		ShopCenter.PUT("/chart", hander.Auth(), controller.Update)
		//order 生成订单
		ShopCenter.POST("/order", hander.Auth(), controller.MakeOrder)
	}
	//管理组

	admin := r.Group("/admin")
	{

		//登录
		admin.POST("/login", controller.RootLogin)
		//展示所有订单 也写个分类吧
		admin.Use(hander.Auth())
		admin.Use(hander.RootAccess())
		admin.GET("/order", controller.RootAll)
		//更新订单
		admin.PUT("/order", controller.UpdateONeOrder)
		//删除订单
		admin.DELETE("/order", controller.DeleteUserOrder)
		//增加商品
		admin.POST("/goods", controller.AddGoods)
		//获取商品
		admin.GET("/goods", controller.AllShop)
		//更新商品信息
		admin.PUT("/goods", controller.UpdateGoods)
		//删除商品
		admin.DELETE("/goods", controller.DeleteGoods)
		//添加订阅分类
		admin.POST("/category", controller.AddCategory)
		//所有分类
		admin.GET("/category", controller.AllCategory)
		//登出
		admin.GET("/logout", controller.RootLogout)
	}
	r.Run(":8080")
	//runtls 实现https 访问 用的是腾讯的ssl 证书 不存在爆红
	//r.RunTLS(":443", "test.pem", "test.key")

}
