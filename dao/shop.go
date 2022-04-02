package dao

import (
	"JD/models"
	"errors"
	"gorm.io/gorm"
)

func AllShops() *models.AllGoos {
	var (
		goods []models.OrmGoodsLIst
	)

	var all models.AllGoos
	//查询商品id以便于 在另外一张表用
	DB.Find(&goods)
	// 嘿嘿 开始蹩脚并发

	for k, v := range goods {
		DB.Where(&models.OrmGoodInfo{
			Gid: v.Gid,
		}).Find(goods[k])
	}
	all.Goods = goods
	//蹩脚 并发

	return &all

}

func MakeOrder(order models.Order) (bool, string) {

	//检验数据是否合法
	var uid []int
	if result := DB.Where(order.ChartId).Find(&uid); result.Error != nil {
		return false, "订单信息错误"
	}
	for _, v := range uid {
		if v != order.Uid {
			return false, "订单信息错误"
		}
	}

	tx := DB.Begin()
	var (
		orders []models.OrmShopChart
		gid    []int
	)
	//插入用户订单 进行订单创建
	tx.Where("chart_id in (?)", order.ChartId).Find(&orders)
	for k, v := range orders {
		if result := tx.Create(&models.OrmUserOrder{
			Uid:   v.Uid,
			Gid:   v.Gid,
			Count: v.Count,
		}); result.Error != nil {
			tx.Rollback()
			return false, "订单创建失败"
		}
		gid[k] = v.Gid

	}
	//对 用户购物车进行删除相关订单
	result := tx.Delete(models.OrmShopChart{}, "chart_id in (?)", order.ChartId)
	if result.Error != nil {
		tx.Rollback()
		return false, "订单创建失败"
	}
	var (
		balance float64
		Price   []float64
		all     float64
	)
	//查询余额是否满足要求
	result = tx.Select("balance").Where(&models.OrmUserInfo{
		Uid: order.Uid,
	}).Scan(&balance)
	if result.Error != nil {
		tx.Rollback()
		return false, "订单创建失败"
	}
	//查询商品价格
	tx.Select("price").Where("Gid in (?)", gid).Scan(Price)
	for _, v := range Price {
		all += v
	}
	if balance < all {
		tx.Rollback()
		return false, "你个穷逼"
	}
	//更新用余额
	if result = tx.Model(&models.OrmUserInfo{}).Update("balance", balance-all); result.Error != nil {
		tx.Rollback()
		return false, "订单创建失败"
	}
	//对商品销量进行增加
	result = tx.Model(&models.OrmGoodInfo{}).Where("GId in (?)", gid).Update("sales", gorm.Expr("sales+1"))
	if result.Error != nil {
		tx.Rollback()
		return false, "订单创建失败"
	}
	//提交事务
	tx.Commit()
	return true, "订单提交成功"
}
func AllOrder(user models.User) (bool, *models.UserOrder) {
	var allorder []models.OrmUserOrder
	DB.Where("uid=?", user.Uid).Find(&allorder)
	var all models.UserOrder
	all.Allorder = allorder
	return true, &all
}
func UpdateOrder(order models.UpdateOrder) error {

	var Orderstate string
	DB.Where(
		&models.OrmUserOrder{Oid: order.Oid}).Select("state").Find(&Orderstate)
	if Orderstate != "已支付" {
		if Orderstate == "已完成" {
			err := errors.New("订单已经完成了哦")
			return err
		}
		err := errors.New("还未发货 着啥急")
		return err
	}
	if result := DB.Model(
		&models.OrmUserOrder{}).Where(
		"oid=?", order.Oid).Update(
		"state", "已完成"); result.Error != nil {
		err := errors.New("确认收获失败")
		return err
	}

	return nil

}

func DeleteOrder(order models.UpdateOrder) error {
	var info models.OrmUserOrder
	result := DB.Table("user_order").Where("oid=?").Scan(&info)
	if result.Error != nil || result.RowsAffected == 0 {
		err := errors.New("订单信息不存在")
		return err
	}

	if info.Uid != order.Uid {
		err := errors.New("怎么可以动别人的订单呢")
		return err
	}
	if info.State != "已完成" {
		err := errors.New("客官 订单还没跑完呢 不要抛弃我呀")
		return err
	}
	DB.Delete(models.OrmUserOrder{}, "oid=?", order.Oid)
	return nil
}
func Commit(commit models.Commit) (bool, string) {
	var order models.OrmUserOrder
	DB.Where("oid=?").Find(&order)

	if order.Uid != commit.Uid || order.State != "已完成" {
		return false, "商品状态错误"
	}
	var com models.OrmCommit
	result := DB.Where("oid=?").Find(&com)

	if result.RowsAffected != 0 {
		return false, "评论提交失败 已经评论过了"
	}
	var onecommit models.OrmCommit
	onecommit.Commit = commit.Commit
	onecommit.Gid = com.Gid
	onecommit.Oid = commit.Oid
	if result := DB.Create(&onecommit); result.Error != nil {
		return false, "评论提交失败"
	}
	result = DB.Model(
		&models.OrmGoodInfo{}).Where(
		"GId=?", order.Gid).Update(
		"commit", gorm.Expr("commit+1"))
	if result.Error != nil {
		return false, "评论提交失败"
	}

	return true, "评论提价成功"
}
func GetCommit(commit models.Commits) *models.AllCommit {

	var info models.OrmGoodInfo
	DB.Where("GId=?", commit.Gid).Find(&info)
	var allcommit []models.OrmCommit
	DB.Where(&models.OrmCommit{Gid: info.Gid}).Find(&allcommit)
	var all models.AllCommit
	all.Gid = info.Gid
	all.Introduce = info.Introduce
	all.Onecomit = allcommit
	return &all
}
func Class() (*models.AllShop, error) {
	var t []models.OrmClass
	var all models.AllShop
	result := DB.Table("goods_class").Select("type").Find(&t)
	if result.Error != nil {
		err := errors.New("查询失败")
		return nil, err
	}
	for _, v := range t {
		var onelist []int
		DB.Table("goods_list").Where("type=?", v).Select("Gid").Find(&onelist)
		for _, m := range onelist {
			var list models.OrmGoodsLIst
			DB.Table("goods_list").Where("Gid = ?", m).Find(&list)
			var onegood models.OrmGoodInfo
			DB.Table("goods_info").Where("Gid = ?", m).Find(&onegood)
			var one models.OrmOneGood
			one.OrmGoodsLIst = list
			one.OrmGoodInfo = onegood
			v.Goods = append(v.Goods, one)
		}
	}
	all.All = t
	return &all, nil

}
