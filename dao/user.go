package dao

import (
	"JD/models"
	"errors"
	"gorm.io/gorm"
	"os"
)

func Register(name string, _ string, number string) error {
	var U models.Register
	//查看用户是否存在
	result := DB.Table("user_info").Where(&models.Register{
		Username: name,
	}).Or(&models.Register{
		Number: number,
	}).Find(&models.Register{})
	if result.RowsAffected != 0 {
		err := errors.New("用户已存在 请尝试登录")
		return err
	}
	if result.Error != nil {
		err := errors.New("注册失败 待会儿再试试吧")
		return err
	}
	//执行注册
	result = DB.Table("user_info").Create(&U)
	if result.Error != nil {
		err := errors.New("注册失败 待会儿再试试吧")
		return err
	}

	return nil

}

type User models.Register

func (User) TableName() string {
	return "user_info"
}
func Find(user models.Register) error {
	var temple User
	result := DB.Where(&User{
		Username: user.Username,
	}).Find(&temple)
	// 进行检测是否傻逼操作

	if result.Error != nil || result.RowsAffected == 0 {
		err := errors.New("找回密码失败 用户信息不存在")
		return err
	}
	if temple.Number != user.Number {
		err := errors.New("手机号不匹配，再试试把")
		return err
	}
	if temple.Password == user.Password {
		err := errors.New("密码不能和原来的相同哦")
		return err
	}
	// 更新密码

	DB.Model(&models.Register{}).Where(&models.Register{
		Username: user.Username,
	}).Updates(user.Password)

	return nil

}

func Login(u models.Login) (*models.BasicInfo, error) {
	var info models.UserInfo
	if result := DB.Table("user_info").Where(&models.Register{Username: u.Username}).Find(&info); result.Error != nil {
		if result.RowsAffected == 0 {
			err := errors.New("你还没注册登录个屁")
			return nil, err
		}
		err := errors.New("登录失败")
		return nil, err
	}
	if info.Word != u.Password {
		err := errors.New("密码错误")
		return nil, err

	}
	//返回
	var basicinfo models.BasicInfo
	basicinfo.Uid = info.Uid
	basicinfo.Username = info.Name
	return &basicinfo, nil
}
func AdminLogin(name string, word string) (bool, string) {
	type admin struct {
		name string
		word string
	}
	var admins admin
	if result := DB.Table("admin_info").Where("name=?", name).Find(&admins); result.Error != nil || result.RowsAffected == 0 {
		if result.RowsAffected == 0 {
			return false, "用户不存在"
		}
		return false, "登录失败"
	}
	if admins.word != word {
		return false, "密码错误"
	}
	return true, "密码正确"

}

type Userinfo models.UserInfo

// TableName 修改表明
func (Userinfo) TableName() string {
	return "user_info"
}
func SaveFile(url string, user models.BasicInfo) (string, error) {
	var info Userinfo
	DB.Where(&Userinfo{
		Name: user.Username,
		Uid:  user.Uid,
	}).Find(&info)
	//查询用户头像是否是默认头像
	if info.Image == "default.png" {

	} else {
		err := os.Remove("./static/" + info.Image)
		if err != nil {

		}
	}
	result := DB.Where(&Userinfo{
		Name: user.Username,
		Uid:  user.Uid,
	}).Updates(&Userinfo{
		Image: url,
	})

	if result.Error != nil {
		err := errors.New("文件上传失败")
		return "", err
	}
	return "文件上传成功", nil
}

func GetBalance(u models.BasicInfo) (*float64, error) {

	var info Userinfo
	result := DB.Where(&Userinfo{
		Name: u.Username,
		Uid:  u.Uid,
	}).Find(&info)
	if result.Error != nil || result.RowsAffected == 0 {
		err := errors.New("查询失败 无匹配信息")
		return nil, err
	}
	return &info.Balance, nil

}
func ChargeBalance(u models.Balance) (bool, string) {
	//充值

	result := DB.Where(&Userinfo{
		Uid:  u.Uid,
		Name: u.Username,
	}).Update("balance", gorm.Expr("balance+?", u.Balance))
	if result.RowsAffected == 0 || result.Error != nil {
		return false, "充值失败"
	}
	return true, "充值成功"
}

func AddChart(chart models.AddChart) (string, error) {
	var all []models.ChartShop
	//查看是否重复
	result := DB.Table("shop_chart").Where("uid=? and gid=?", chart.Uid, chart.Gid).Find(&all)
	if result.RowsAffected != 0 {
		err := errors.New("你已经加入老人购物车 想干啥")
		return "", err
	}
	//插入数据
	result = DB.Table("shop_chart").Create(&models.AddChart{
		Uid:   chart.BasicInfo.Uid,
		Gid:   chart.Gid,
		Count: chart.Count,
	})
	if result.RowsAffected == 0 {
		err := errors.New("加入购物车失败")
		return "", err
	}
	return "你的宝贝已经躺在购物车里了哦", nil

}
func AllChart(user models.Userinfo) (*models.AllChart, error) {
	all := models.AllChart{}
	var T []models.ChartShop
	result := DB.Table("shop_chart").Where("uid=?", user.Uid).Find(&T)
	all.BasicInfo = user.BasicInfo
	if result.Error != nil {
		err := errors.New("查询失败")
		return nil, err
	}
	// 本来准备丢sql循环 这样不好的 最后还是丢进去了 可能因为数据库结构不合理吧

	for k, v := range T {
		DB.Table("goods_list").Select("Name").Where("Gid =?", v.Gid).Scan(&T[k].Good)
	}
	all.ChartList = T
	return &all, nil

}
func UpdateChart(chart models.ShopChart) (bool, string) {

	//将Count设置为0 就意味着删除
	if chart.Count == 0 {

		//查看数据是否合法
		var OneChart int
		result := DB.Table("shop_chart").Where("chart_id=?", chart.ChartId).Select("uid").Scan(&OneChart)
		if result.Error != nil {
			return false, "失败 呜呜呜"
		}
		if OneChart != chart.Uid {
			return false, "你怎么能动别人的订单呢"
		}
		result = DB.Table("shop_chart").Where("chart_id=?", chart.ChartId).Delete(models.ChartShop{})
		if result.Error != nil {
			return false, "失败 呜呜  呜"

		}

		return true, "宝贝忍痛离开了购物车"

	}
	var Temple models.ShopChart
	result := DB.Table("shop_chart").Where("chart_id=?", chart.ChartId).Find(&Temple)
	if result.Error != nil || result.RowsAffected == 0 {
		return false, "失败  没有查到该订单"
	}
	if Temple.Uid != chart.Uid {
		return false, "你怎么能动别人的订单呢"
	}
	result = DB.Table("shop_chart").Update("Count=?", chart.Count)
	if result.Error != nil {
		return false, "失败 怎那么就失败呢"
	}
	return true, "操作成功"
}

func MyInfo(User models.BasicInfo) (*models.MyInfo, error) {

	var UserInfo models.MyInfo
	UserInfo.BasicInfo = User
	//查询用户基本信息
	DB.Table("user_info").Where(&models.UserInfo{
		Uid: User.Uid,
	}).Select("balance,image").Scan(&UserInfo)
	var link = "https://sanser.ltd/static/"

	UserInfo.ImageUrl += link
	var Category []models.Category
	//查询用户订单
	DB.Table("order_state_type").Find(&Category)
	UserInfo.Category = Category
	var temple models.AllOrder
	//对订单进行归入
	for key, value := range UserInfo.Category {
		DB.Table("user_order").Where("uid=? and state=?", User.Uid, value.State).Find(&temple)
		UserInfo.Category[key].Order = append(UserInfo.Category[key].Order, temple)
	}
	UserInfo.BasicInfo = User
	return &UserInfo, nil
}
func IdInfo(id int) (models.IdToken, error) {
	var idtoken models.IdToken
	DB.Where("uid=?", id).Find(&idtoken)
	return idtoken, nil
}
