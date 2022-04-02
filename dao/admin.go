package dao

import (
	"JD/models"
	"JD/utils"
	"os"
	"strconv"
	"sync"
	"time"

	"errors"

	"github.com/gin-gonic/gin"
)

func GetAllOrder() (*models.AllInfo, error) {
	var AllInfo models.AllInfo
	var types []string
	result := DB.Table("Order_state_type").Find(&types)
	if result.Error != nil {
		err := errors.New("数据查询失败")
		return nil, err
	}

	for _, v := range types {
		var all models.All
		all.State = v
		var one []models.OneOrder
		DB.Table("user_order").Where("state=?", v).Find(&one)
		all.Order = one
		AllInfo.All = append(AllInfo.All, all)

	}

	return &AllInfo, nil
}
func OrderChange(order models.UpdateUserOrder) (string, error) {
	var uid int
	result := DB.Table("user_order").Where("oid=?", order.Oid).Select("uid").Find(&uid)
	if result.Error != nil {
		err := errors.New("订单信息错误")
		return "", err
	}
	utils.AddCategory(strconv.Itoa(order.Oid))
	utils.Subscribe(uid, strconv.Itoa(order.Oid))
	result = DB.Table("user_order").Where("oid=?", order.Oid).Update("state", order.State)
	if result.Error != nil {
		err := errors.New("订单信息错误")
		return "", err
	}
	t := time.Now()
	utils.AddNewInfo(strconv.Itoa(order.Oid), "订单状态变化", "您的订单状态发生变化", t)
	return "操作成功", nil
}
func DeleteUserOrder(update models.UpdateUserOrder) (string, error) {
	var state string
	DB.Table("user_order").Select("state").Where("oid=?").Scan(&state)
	if state != "已完成" {
		err := errors.New("订单信息错误")
		return "", err
	}
	result := DB.Delete(models.OrmUserOrder{}, "oid=?", update.Oid)
	if result.Error != nil {
		err := errors.New("订单信息删除错误")
		return "", err
	}
	return "", nil

}
func AddGoods(goods models.GoodsAdd, c *gin.Context) (string, error) {
	tx := DB.Begin()
	url, err := utils.SaveFile(goods.Image, c)
	if err != nil {

		//err = errors.New("商品添添加失败 请重试")
		tx.Rollback()
		return "", err
	}
	var onegood models.OrmGoodsLIst
	onegood.Type = goods.Category
	onegood.Url = url
	onegood.Name = goods.Gname
	result := tx.Create(&onegood)
	if result.Error != nil {
		tx.Rollback()
		err := errors.New("添加失败")
		return "", err
	}
	var one models.OrmGoodInfo
	one.Introduce = goods.Introduce
	one.Price = float64(int(goods.Price))
	if result := tx.Create(&one); result.Error != nil {
		err := errors.New("商品添加失败")
		tx.Rollback()
		return "", err
	}
	tx.Commit()
	return "商品添加成功", nil

}
func UpdateGoods(goods models.UpdateGoods, c *gin.Context) (string, error) {
	s := sync.WaitGroup{}
	s.Add(1)
	go func() {
		t := time.Now()
		utils.AddNewInfo(strconv.Itoa(goods.Gid), "商品信息修改", "您关注的商品信息发生变动", t)
		s.Done()
	}()
	tx := DB.Begin()

	if goods.Image != nil {
		var url string
		result := DB.Table("goods_list").Select("url").Find(&url)
		if result.Error != nil {

			err := errors.New("数据更新失败 请稍后重试")
			return "", err
		}
		err := os.Remove("/www/static/" + url)
		if err != nil {
			err = errors.New("数据更新失败 请稍后重试")
			return "", err
		}
		url, err = utils.SaveFile(goods.Image, c)
		if err != nil {
			tx.Rollback()
			err = errors.New("数据更新失败 请稍后重试")
			return "", err
		}
		result = DB.Model(&models.OrmGoodsLIst{}).Where("Gid=?", goods.Gid).Update("url", url)
		if result.Error != nil {
			err = errors.New("数据更新失败 请稍后重试")
			tx.Rollback()
			return "", err
		}
		s.Wait()
	}
	if goods.Gname != "" {
		result := DB.Model(&models.OrmGoodsLIst{}).Where("Gid=?", goods.Gid).Update("name", goods.Gname)
		if result.Error != nil {
			err := errors.New("数据更新失败 请稍后重试")
			tx.Rollback()
			return "", err
		}
	}
	if goods.Introduce != "" {
		result := DB.Model(&models.OrmGoodsLIst{}).Where("Gid=?", goods.Gid).Update("introduce", goods.Introduce)
		if result.Error != nil {
			err := errors.New("数据更新失败 请稍后重试")
			tx.Rollback()
			return "", err
		}
	}
	if goods.Price != 0 {
		result := DB.Model(&models.OrmGoodsLIst{}).Where("Gid=?", goods.Gid).Update("price", &goods.Price)
		if result.Error != nil {
			err := errors.New("数据更新失败 请稍后重试")
			tx.Rollback()
			return "", err
		}
	}
	tx.Commit()
	return "小商品修改成功了哦", nil
}

func DeleteGoods(gid string) (string, error) {
	tx := DB.Begin()

	//var wait sync.WaitGroup
	result := DB.Delete(models.OrmCommit{}, "Gid=?", gid)
	if result.Error != nil {
		err := errors.New("商品删除失败 再试试吧")
		tx.Rollback()
		return "", err
	}
	result = DB.Delete(models.OrmShopChart{}, "gid=？", gid)
	if result.Error != nil {
		err := errors.New("商品删除失败 再试试吧")
		tx.Rollback()
		return "", err
	}
	result = DB.Delete(models.OrmGoodInfo{}, "Gid=?", gid)
	if result.Error != nil {
		err := errors.New("商品删除失败 再试试吧")
		tx.Rollback()
		return "", err
	}
	var url string
	DB.Table("goods_list").Select("url").Where("Gid=?", gid).Scan(&url)

	err := os.Remove("/www/static/" + url)
	if err != nil {
		err = errors.New("商品删除失败 再试试吧")
		tx.Rollback()
		return "", err
	}
	result = DB.Delete(models.OrmGoodsLIst{}, "Gid=?", gid)
	if result.Error != nil {
		err = errors.New("商品删除失败 再试试吧")
		tx.Rollback()
		return "", err
	}

	tx.Commit()
	return "商品删除成功", nil

}
